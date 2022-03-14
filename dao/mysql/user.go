package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"webapp/models"
)

// 把数据库操作封装成函数
// logic层按需调用

const secret = "secret"
var (
	ErrorUserExist = errors.New("用户已存在")
	ErrorUserNotExist = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

func CheckUesrExist(username string) (err error) {
	// 执行sql语句
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}

	if count > 0 {
		return ErrorUserExist
	}

	return
}

func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)

	// 执行sql语句入库
	sqlStr := "insert into user(user_id, username, password) values(?, ?, ?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)

	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error){
	oPassword := user.Password
	sqlStr := "select user_id, username, password from user where username = ?"

	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	if user.Password != encryptPassword(oPassword) {
		return ErrorInvalidPassword
	}
	
	return 
}

func GetUserByID(uid int64) (username string, err error) {
	sqlStr := "select username from user where user_id = ?"
	err = db.Get(&username, sqlStr, uid)
	return
}	