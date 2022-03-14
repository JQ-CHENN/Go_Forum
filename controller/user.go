package controller

import (
	"errors"
	"webapp/dao/mysql"
	"webapp/logic"
	"webapp/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamSignUp) 
	if err := c.ShouldBindJSON(p); err != nil {
		//参数有误
		zap.L().Error("SignUp with invild parmas", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}
	// 对请求参数进行校验

	
	// 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic signup failed", zap.Error(err))

		if errors.Is(err, mysql.ErrorUserExist){
			ResponseError(c, CodeUserExit)
			return
		}

		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 将获取的参数与数据库中的参数进行校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invild parmas", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}
	
	// 登录处理
	token, err := logic.Login(p); 
	if err != nil {
		zap.L().Error("logic Login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExit)
			return
		} 
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 返回响应
	ResponseSuccess(c, token)
}

