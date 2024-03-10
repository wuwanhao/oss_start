package router

import (
	"api_service/api/controller"
	"api_service/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {

	router := gin.New()
	// 宕机时自动恢复
	router.Use(gin.Recovery())
	// 跨域中间件
	router.Use(middleware.Cors())
	// 注册路由
	register(router)
	return router
}

// register 路由接口
func register(router *gin.Engine) {
	router.POST("/file", controller.PutObject)
}
