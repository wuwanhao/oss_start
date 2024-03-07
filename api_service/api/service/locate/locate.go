package locate

import (
	"connector/rabbitmq"
	"strconv"
	"time"
)

// 从服务节点定位文件
func Locate(name string) string {

	q := rabbitmq.New()
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
