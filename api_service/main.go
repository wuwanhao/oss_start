package main

import (
	"api_service/api/service/heartbeat"
	"api_service/common/config"
	"api_service/router"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// 接收来自数据服务节点的心跳
	go heartbeat.ListenHeartbeat()

	// 设置启动模式
	gin.SetMode(config.Config.Server.Mode)
	// 初始化路由
	routers := router.InitRouter()
	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: routers,
	}
	// 启动服务
	log.Println("Starting Server...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen failed: %v\n", config.Config.Server.Address)
		}
	}()
	log.Printf("success, listen: %v\n", config.Config.Server.Address)

	// 监听退出消息
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exit!")
}
