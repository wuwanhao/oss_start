package main

import (
	"api_service/api/service/heartbeat"
	"api_service/common/config"
	"api_service/router"
	"common_service/logs"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// 初始化日志库
	logs.InitLog(config.Config.Server.Name)
	// 设置启动模式
	gin.SetMode(config.Config.Server.Mode)
	// 接收来自数据服务节点的心跳
	go heartbeat.ListenHeartbeat()
	// 初始化路由
	routers := router.InitRouter()
	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: routers,
	}
	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error("server listen failed: %v\n", config.Config.Server.Address)
		}
	}()
	logs.Info("server start success, listen: %v\n", config.Config.Server.Address)

	// 监听退出消息
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error("server Shutdown: %v\n", err)
	}
	logs.Info("server exit!")
}
