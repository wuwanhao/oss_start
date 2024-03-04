package objects

import (
	"api_service/heartbeat"
	"api_service/locate"
	objectStream "api_service/objectStream"
	"fmt"
	"io"
	"log"
	"net/http"
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

	w.WriteHeader(http.StatusMethodNotAllowed)
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
	// 1.拿到要上传的文件名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 2.上传文件
	c, err := storeObject(object, r.Body)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(c)
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
