package redis

const (
	KeyPrefix = "webapp"
	KeyPostTimeZSet = "post:time" // 按发帖时间 score:时间 member: 帖子id
	KeyPostScoreZset = "post:score" // 按投票分数 score: 分数 member: 帖子id
	KeyPostVoteZsetPre = "post:voted:" // 记录用户及投票类型，参数是post id，member是用户id
	KeyCommunitySetPre = "community:" // 保存每个社区下帖子的id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}