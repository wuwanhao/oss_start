package heartbeat

import (
	"data_service/utils/rabbitmq"
	"os"
	"time"
)

// 发送心跳
func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	for {
		q.Publish("apiServers", os.Getenv("LISTENLADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
