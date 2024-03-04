package locate

import (
	"connector/rabbitmq"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 用于向数据服务节点群发定位消息并接收反馈
func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取get的key
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// 从服务节点定位文件
func Locate(name string) string {

	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	// 向名为dataServers的exchange群发定位信息
	q.Publish("dataServers", name)
	c := q.Consume()

	// 起一个协程，1S后定位不到，关闭MQ连接
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()

	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
