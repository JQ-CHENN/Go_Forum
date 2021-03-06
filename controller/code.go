package controller

type ResCode int64
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParma
	CodeUserExit
	CodeUserNotExit
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
)

var codeMsgMap = map[ResCode]string {
	CodeSuccess: "请求成功",
	CodeInvalidParma: "请求参数错误",
	CodeUserExit: "用户名已存在",
	CodeUserNotExit: "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy: "服务繁忙",
	CodeInvalidToken: "无效的token",
	CodeNeedLogin: "需要登录",
}

func (c ResCode) GetMsg() string {
	msg, ok := codeMsgMap[c]
	if ! ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}