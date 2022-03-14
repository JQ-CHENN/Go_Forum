package controller

import (
	"webapp/logic"
	"webapp/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, CodeInvalidParma)
		return
	}

	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("get current user err", zap.Error(err))
		return
	}
	if err := logic.PostVote(userID, p); err != nil {
		zap.L().Error("Post vote err", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}