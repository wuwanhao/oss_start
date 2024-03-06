package controller

import (
	"api_service/api/service/locate"
	"api_service/common/result"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

// 定位文件
func LocateFile(c *gin.Context) {
	// 获取get的key
	url := c.Request.URL
	log.Println(url)
	info := locate.Locate(strings.Split(url.EscapedPath(), "/")[2])
	result.Success(c, info)
}
