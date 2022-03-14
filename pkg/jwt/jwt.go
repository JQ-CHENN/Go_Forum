package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 24 // 过期时间
var MySecret = []byte("夏天夏天悄悄过去")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID int64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	c := &MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer: "webapp",
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}