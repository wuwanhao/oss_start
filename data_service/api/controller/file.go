package controller

import (
	"data_service/api/service/file"
	"data_service/common/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

// 获取文件流
func GetFile(c *gin.Context) {

	fileService := file.NewFileService(config.Config.Oss.StorageRoot, config.Config.Oss.StorageIndex)
	filename := c.Param("filename")
	getFile, err := fileService.GetFile(filename)
	if err != nil {
		log.Println("[data_service] get file error: ", err)
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Writer, getFile)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
}
