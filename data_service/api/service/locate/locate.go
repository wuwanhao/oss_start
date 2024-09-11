package locate

import (
	"common_service/logs"
	"connector_service/rabbitmq"
	"data_service/common/config"
	"data_service/common/utils"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)

}

// 监听api_service需要定位的文件，返回定位结果
func StartLocate() {
	q := rabbitmq.New()
	defer q.Close()
	// 监控dataServers的消息
	q.Bind(config.Config.RabbitMq.DataExchange)
	c := q.Consume()

	// 遍历这些消息，若在自己的服务中定位到该文件，则向队列中返回消息，带着自己的服务地址
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			logs.Warn("locate msg resolve error: %v\n", e)
			continue
		}
		if Locate(config.Config.Oss.StorageRoot + config.Config.Oss.StorageIndex + object) {
			q.Send(msg.ReplyTo, utils.GetServerIp())
		}
	}
}
