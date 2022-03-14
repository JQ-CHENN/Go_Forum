package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webapp/dao/mysql"
	"webapp/dao/redis"
	"webapp/logger"
	"webapp/pkg/snowflake"
	"webapp/routers"
	"webapp/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 加载配置(文件)
	if err := settings.Init(); err != nil {
		fmt.Println("init settings err: ", err)
		return
	}
	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("init log err: ", err)
		return
	}
	defer zap.L().Sync()
	// 初始化mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Println("init mysql err: ", err)
		return
	}
	defer mysql.Close()
	// 初始化redis连接
	if err := redis.Init(); err != nil {
		fmt.Println("init redis err: ", err)
		return
	}
	defer redis.Close()

	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machineID")); err != nil {
		fmt.Println("init snowfake err: ", err)
		return
	}

	// 注册路由
	r := routers.Setup()
	// 启动服务
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func () {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	// 创建一个5秒超时的context
	ctx, cancel:= context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}