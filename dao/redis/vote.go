package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSec = 7 * 24 * 3600
	scorePerVote = 432 // 每一票分数
)
var (
	ErrorVoteTimeExpire = errors.New("投票已过期")
	ErrorVoteRepeat = errors.New("不允许重复投票")
)

func CreatePost(id int64) (err error) {
	pipeline := client.TxPipeline() // 事务处理
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: id,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: id,
	})

	_, err = pipeline.Exec()
	return 
}

// 投票功能
// 用户投票的数据（赞成，反对）

// 投一票加432分

/*
	投票情况
	UpLow = 1    1）之前没有投过票，2）之前投反对票
	UpLow = 0    1）之前投过赞成票 2）之前投反对票
	UpLow = -1   1）之前没有投过票 2）之前投过赞成票
	-->更新分数和投票纪录

	限制： 1)每个帖子发表后一个星期内可以投票，到期后将redis中保存的投票数据存到mysql
		  2)到期删除相应的Key
*/
func PostVote(userID, postID string, UpLow float64) error {
	// 判断投票的限制
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix()) - postTime > oneWeekInSec {
		return ErrorVoteTimeExpire
	}

	// 更新分数
	// 先查之前的投票纪录
	oUpLow := client.ZScore(getRedisKey(KeyPostVoteZsetPre + postID), userID).Val()
	diff := math.Abs(oUpLow - UpLow) // 计算两次投票的差值
	if UpLow == oUpLow {
		return ErrorVoteRepeat
	}

	var op float64
	if UpLow > oUpLow {
		op = 1
	} else {
		op = -1
	}

	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op * diff * scorePerVote, postID).Result()

	// 记录用户给帖子投票的数据 
	if UpLow == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVoteZsetPre + postID), userID)

	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVoteZsetPre + postID), redis.Z{
			Score: UpLow,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}

// 根据id列表查询帖子的赞数
func GetPostVoteUps(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVoteZsetPre + id)

	// 	// 查找key中分数是1的元素的数量
	// 	v := client.ZCount(key, "1", "1").Val()
	// 	data = append(data, v)
	// }

	pipeline := client.Pipeline() // 一次发送多条命令，减少RTT
	for _, id := range ids {
		key := getRedisKey(KeyPostVoteZsetPre + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	} 
	return
}
