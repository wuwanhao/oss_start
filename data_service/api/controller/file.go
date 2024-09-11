package controller

import (
	"common_service/logs"
	"data_service/api/service/file"
	"data_service/common/config"
	"data_service/common/result"
	"github.com/gin-gonic/gin"
	"io"
)

// 获取文件流
func GetFile(c *gin.Context) {

	// 文件服务实例
	fileService := file.NewFileService(config.Config.Oss.StorageRoot, config.Config.Oss.StorageIndex)
	filename := c.Query("file_name")
	getFile, err := fileService.GetFile(filename)
	if err != nil {
		logs.Warn("get file error: %v\n", err)
		result.Failed(c, int(result.ApiCode.FILE_NOT_FOUND), result.ApiCode.GetMessage(result.ApiCode.FILE_NOT_FOUND))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Writer, getFile)
	if err != nil {
		logs.Warn("copy file byte transfer error: %v\n", err)
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
	logs.Info("[data_service] fileName: %v\n", fileName)

	postFile := c.Request.Body
	if postFile == nil {
		logs.Warn("put body is empty")
		result.Failed(c, int(result.ApiCode.FILE_PUT_EMPTY_ERROR), result.ApiCode.GetMessage(result.ApiCode.FILE_PUT_EMPTY_ERROR))
		return
	}
	defer postFile.Close()

	// 2.执行文件保存
	fileService := file.NewFileService(config.Config.Oss.StorageRoot, config.Config.Oss.StorageIndex)
	err := fileService.PutFile(postFile, fileName)
	if err != nil {
		logs.Warn("put file error: %v\n", err)
		result.Failed(c, int(result.ApiCode.FILE_UPLOAD_ERROR), result.ApiCode.GetMessage(result.ApiCode.FILE_UPLOAD_ERROR))
		return
	}

	result.Success(c, fileName)
	return
}
