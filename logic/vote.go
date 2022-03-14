package logic

import (
	"strconv"
	"webapp/models"
	"webapp/dao/redis"
)



func PostVote(userID int64, p *models.ParamVoteData) error {
	return redis.PostVote(strconv.Itoa(int(userID)), strconv.Itoa(int(p.PostID)), float64(p.UpLow))
}