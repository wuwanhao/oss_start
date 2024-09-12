package objects

import (
	"api_service/api/service/heartbeat"
	"api_service/api/service/locate"
	"api_service/api/service/objectStream"
	"api_service/common/result"
	"api_service/common/utils"
	"common_service/logs"
	"connector_service/es"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 删除 todo
func del(w http.ResponseWriter, r *http.Request) {

	// 获取文件名
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, err := es.SearchLatestVersion(name)
	if err != nil {
		logs.Warn("Err search latest version: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// put一个同名，版本加一，但是大小为0，hash为空字符串的元数据，表示这是一个删除标记
	e := es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		logs.Warn("Err put metadata: %v\n", e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// todo 接口服务的get方法
func getFile(c *gin.Context, fileName string, version int) {

	// 1.根据文件名获取对象的元数据信息
	meta, err := es.GetMetadata(fileName, version)
	if err != nil {
		logs.Warn("get object metadata error: ", err)
		result.Failed(c, int(result.ApiCode.ERROR_GET_METADATA), result.ApiCode.GetMessage(result.ApiCode.ERROR_GET_METADATA))
		return
	}

	// 当前要查找的对象已经被删除
	if meta.Hash == "" {
		result.Failed(c, int(result.ApiCode.ERROR_OBJECT_DELETED), result.ApiCode.GetMessage(result.ApiCode.ERROR_OBJECT_DELETED))
		return
	}

	// 4.因为之前是用元数据的hash作为作为对象在数据服务中对应的name，所以此处通过hash去作为对象的name
	object := url.PathEscape(meta.Hash)

	// 5.从数据服务中拿到对象的原始数据
	fileStream, err := getStream(object)
	if err != nil {
		logs.Warn("get object stream from dataService error: %v", err)
		result.Failed(c, int(result.ApiCode.ERROR_GET_FILE), result.ApiCode.GetMessage(result.ApiCode.ERROR_GET_FILE))
		return
	}
	_, err = io.Copy(c.Writer, fileStream)
	return

}

// 从数据服务中拿到对象的原始数据
func getStream(object string) (io.Reader, error) {

	// 1.获取该文件存在的数据服务节点
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate failed", object)
	}
	// 2.向该数据服务节点请求数据
	return objectStream.NewGetStream(server, object)
}

// 文件上传
func UploadObject(c *gin.Context, hash, name string) error {

	// 1.将文件上传到数据服务
	statusCode, err := storeObject(url.PathEscape(hash), c.Request.Body)
	if err != nil {
		logs.Warn("storage object error: %v", err)
		return err
	}
	if statusCode != http.StatusOK {
		return err
	}

	// 2.拿到文件大小
	size := utils.GetSizeFromHeader(c.Request.Header)

	// 3.存储文件信息到ES
	uploadResult := es.AddVersion(name, hash, size)
	if uploadResult != nil {
		logs.Warn("es add file version info error: %v", uploadResult)
		return err
	}

	return nil
}

// 存储文件
func storeObject(object string, r io.ReadCloser) (int, error) {
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	// io.Copy会调用stream里面的Write方法，将r的内容写入Stream
	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}

func putStream(object string) (*objectStream.PutStream, error) {
	// 随机选择一个数据服务
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("Can not find any data server")
	}
	return objectStream.NewPutStream(server, object), nil

}
