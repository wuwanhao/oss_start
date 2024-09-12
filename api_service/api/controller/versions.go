package controller

import (
	"common_service/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetAllVersions 获取指定文件的所有版本
func GetAllVersions(c *gin.Context) {
	// 参数检查
	name := c.Param("name")
	if name == "" {
		logs.Warn("missing object name, %v\n", c.Request.URL)
		c.AbortWithStatus(http.StatusBadRequest)
	}

}

func getParams(c *gin.Context) (string, int, int) {
	// 参数检查
	name := c.Param("name")
	if name == "" {
		logs.Warn("missing object name, %v\n", c.Request.URL)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	// 尝试将字符串转换为 int16 和 int
	fromStr := c.DefaultQuery("from", "0")
	sizeStr := c.DefaultQuery("size", "1000")
	var from int
	var size int
	var err error
	if from, err = strconv.Atoi(fromStr); err != nil {
		logs.Warn("invalid 'from' parameter: %v\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if size, err = strconv.Atoi(sizeStr); err != nil {
		logs.Warn("invalid 'size' parameter: %v\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	return name, from, size

}
