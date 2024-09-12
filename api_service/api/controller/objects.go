package controller

import (
	"api_service/api/service/objects"
	"api_service/common/utils"
	"common_service/logs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// 上传文件
func PutObject(c *gin.Context) {

	// 1.获取文件hash
	hash := utils.GetHashFromHeader(c.Request.Header)
	if hash == "" {
		logs.Warn("missing object hash in digest header, %v\n", c.Request.URL)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 2.文件名检查
	fileName := strings.TrimSpace(c.Query("file_name"))
	log.Println(fileName)
	if fileName == "" {
		logs.Warn("File name is empty, %v\n", c.Request.URL)
		c.Status(http.StatusBadRequest)
		return
	}

	// 3.调用服务层处理上传逻辑
	err := objects.UploadObject(c, hash, fileName)
	if err != nil {
		logs.Warn("object upload err: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// 下载文件
func GetObject(c *gin.Context) {
	// GET /objects/<object_name>?version=<version_id›
}

// 删除文件
func DeleteObject(c *gin.Context) {
	// GET /objects/<object_name>?version=<version_id›
}
