package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 错误码
	"msg": 提示信息
	"data" 数据
}
*/
type ResponseData struct {
	Code ResCode	`json:"code"`
	Msg interface{} `json:"msg"`
	Data interface{}`json:"data,omitempty"`
}

// 响应错误消息
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg: code.GetMsg(),
		Data: nil,
	})
}	

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg: msg,
		Data: nil,
	})
}	

// 响应成功
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg: CodeSuccess.GetMsg(),
		Data: data,
	})
}

