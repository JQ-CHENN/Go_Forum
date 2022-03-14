package logic

import (
	"webapp/dao/mysql"
	"webapp/models"
	"webapp/pkg/jwt"
	"webapp/pkg/snowflake"
)

// 业务逻辑处理

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断user是否存在
	if err = mysql.CheckUesrExist(p.Username); err != nil { // 数据库查询出错
		return
	}
	// 生成UID
	userID := snowflake.GenID()

	// 保存进数据库
	user := &models.User{
		UserID: userID,
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	if err := mysql.Login(user); err != nil {
		return "", err
	}

	// 用数据库查询出来的user_id生成token(指针传递)
	return jwt.GenToken(user.UserID, user.Username)
}
