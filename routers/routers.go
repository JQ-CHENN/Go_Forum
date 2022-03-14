package routers

import (
	"net/http"
	"webapp/controller"
	"webapp/logger"
	"webapp/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	v1 := r.Group("/api/v1")
	
	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)

	// 登录业务路由
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		v1.POST("/vote", controller.PostVoteHandler)
		v1.GET("/posts2", controller.GetPostListHandler2)
	}

	// r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	// 如果是登录用户，判断请求头中是否有有效 jwt token
	// 	c.String(http.StatusOK, "pong")
	// })
	
	return r
}