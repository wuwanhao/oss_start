package main

import (
	"context"
	"data_service/common/config"
	"data_service/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	//// 启动协程，向api_service发送心跳
	//go heartbeat.StartHeartbeat()
	//// 启动一个协程，监听dataServers的消息
	//go locate.StartLocate()

	// 设置启动模式
	gin.SetMode(config.Config.Server.Mode)
	// 初始化路由
	routers := router.InitRouter()
	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: routers,
	}
	// 启动服务
	log.Println("[data_service] Starting Server...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[data_service] listen failed: %v\n", config.Config.Server.Address)
		}
	}()
	log.Printf("[data_service] success, listen: %v\n", config.Config.Server.Address)

	// 监听退出消息
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("[data_service] Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("[data_service] Server Shutdown error:", err)
	}
	log.Println("[data_service] Server exit!")
}
