package locate

import (
	"connector_service/rabbitmq"
	"strconv"
	"time"
)

// 从数据节点定位文件，若某一个数据节点定位到该文件，则该节点会返回自身IP地址
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

// 该文件在所有的数据节点中是否存在
func Exist(name string) bool {
	return Locate(name) != ""
}
