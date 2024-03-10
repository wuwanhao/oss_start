package controller

import (
	"data_service/api/service/file"
	"data_service/common/config"
	"data_service/common/result"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

// 获取文件流
func GetFile(c *gin.Context) {

	fileService := file.NewFileService(config.Config.Oss.StorageRoot, config.Config.Oss.StorageIndex)
	filename := c.Query("file_name")
	getFile, err := fileService.GetFile(filename)
	if err != nil {
		log.Println("[data_service] get file error: ", err)
		result.Failed(c, int(result.ApiCode.FILE_NOT_FOUND), result.ApiCode.GetMessage(result.ApiCode.FILE_NOT_FOUND))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Writer, getFile)
	if err != nil {
		log.Println("[data_service] copy file byte transfer error: ", err)
		result.Failed(c, int(result.ApiCode.FILE_BYTE_TRANS_ERROR), result.ApiCode.GetMessage(result.ApiCode.FILE_BYTE_TRANS_ERROR))
		return
	}

	return

}

// 上传文件
func PutFile(c *gin.Context) {

	// 1.获取文件原信息
	fileName := c.Param("file_name")
	if fileName == "" {
		result.Failed(c, int(result.ApiCode.FILE_NAME_CHECK_ERROR), result.ApiCode.GetMessage(result.ApiCode.FILE_NAME_CHECK_ERROR))
		return
	}
	log.Println("[data_service] fileName: ", fileName)

	postFile := c.Request.Body
	if postFile == nil {
		log.Println("[data_service] put body is empty")
		result.Failed(c, int(result.ApiCode.FILE_PUT_EMPTY_ERROR), result.ApiCode.GetMessage(result.ApiCode.FILE_PUT_EMPTY_ERROR))
		return
	}
	defer postFile.Close()

	// 2.执行文件保存
	fileService := file.NewFileService(config.Config.Oss.StorageRoot, config.Config.Oss.StorageIndex)
	err := fileService.PutFile(postFile, fileName)
	if err != nil {
		log.Println("[data_service] put file error: ", err)
		result.Failed(c, int(result.ApiCode.FILE_UPLOAD_ERROR), result.ApiCode.GetMessage(result.ApiCode.FILE_UPLOAD_ERROR))
		return
	}

	result.Success(c, fileName)
	return
}
