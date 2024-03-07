package controller

import (
	"api_service/api/service/locate"
	"api_service/common/result"
	"github.com/gin-gonic/gin"
	"log"
)

// 从数据服务节点定位文件
func LocateFile(c *gin.Context) {

	name := c.Param("name")
	log.Println("name:" + name)
	info := locate.Locate(name)
	result.Success(c, info)
}
