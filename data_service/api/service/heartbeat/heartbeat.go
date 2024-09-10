package heartbeat

import (
	"connector_service/rabbitmq"
	"data_service/utils"
	"time"
)

// 发送心跳，注册到api_service
func StartHeartbeat() {
	q := rabbitmq.New()
	defer q.Close()

	// 5s 一次
	for {
		q.Publish("apiServers", utils.GetServerIp()+":"+utils.GetServerHttpPort())
		time.Sleep(5 * time.Second)
	}
}
