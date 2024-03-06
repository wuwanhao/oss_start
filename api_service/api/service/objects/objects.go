package objects

import (
	"api_service/api/service/heartbeat"
	"api_service/api/service/locate"
	"api_service/api/service/objectStream"
	"api_service/utils"
	"connector/es"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodPut {
		put(w, r)
		return
	}

	if method == http.MethodGet {
		get(w, r)
		return
	}
	if method == http.MethodDelete {
		del(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 删除
func del(w http.ResponseWriter, r *http.Request) {

	// 获取文件名
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, err := es.SearchLatestVersion(name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// put一个同名，版本加一，但是大小为0，hash为空字符串的元数据，表示这是一个删除标记
	e := es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// 接口服务的get方法
func get(w http.ResponseWriter, r *http.Request) {
	// 1.拿到要取得的文件名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := getStream(object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)

}

func getStream(object string) (io.Reader, error) {

	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate failed", object)
	}
	return objectStream.NewGetStream(server, object)
}

// 接口服务的put方法
func put(w http.ResponseWriter, r *http.Request) {
	// 1.获取文件hash
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c, e := storeObject(url.PathEscape(hash), r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	// 2.拿到要上传的文件名
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 3.拿到文件大小
	size := utils.GetSizeFromHeader(r.Header)
	// 4.上传文件
	//uploadResult := es.AddVersion(name, hash, size)
	fmt.Println(name)
	fmt.Println(size)
	fmt.Println(hash)
	// if uploadResult != nil {
	// 	log.Println(uploadResult)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// }

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
