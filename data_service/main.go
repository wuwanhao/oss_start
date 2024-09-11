package main

import (
	"common_service/logs"
	"context"
	"data_service/api/service/heartbeat"
	"data_service/api/service/locate"
	"data_service/common/config"
	"data_service/router"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// 初始化日志库
	logs.InitLog(config.Config.Server.Name)
	// 启动协程，向api_service发送心跳
	go heartbeat.StartHeartbeat()
	// 启动一个协程，监听api_service需要定位的文件，返回定位结果
	go locate.StartLocate()
	// 设置启动模式
	gin.SetMode(config.Config.Server.Mode)
	// 初始化路由
	routers := router.InitRouter()
	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: routers,
	}
	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error("Listen address failed: %v\n", config.Config.Server.Address)
		}
	}()
	logs.Info("Start success, listen: %v\n", config.Config.Server.Address)

	// 监听退出消息
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error("Server Shutdown error: %v\n", err)
	}
	logs.Info("Server already shut down!")
}
