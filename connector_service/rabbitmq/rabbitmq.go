package rabbitmq

import (
	"api_service/common/config"
	"common_service/logs"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
)

// RabbitMQ 对象
type RabbitMQ struct {
	channel  *amqp091.Channel // 连接通道
	Name     string           // 队列名
	exchange string           // 队列绑定的交换机名
}

// 新建一个RabbitMQ对象
func New() *RabbitMQ {
	// rabbitMQ 连接url
	conn_param := "amqp://" + config.Config.RabbitMq.Username + ":" + config.Config.RabbitMq.Password + "@" + config.Config.RabbitMq.Host + ":" + strconv.Itoa(config.Config.RabbitMq.Port) + "/" + config.Config.RabbitMq.Virtualhost
	conn, err := amqp091.Dial(conn_param)
	if err != nil {
		logs.Error("Connect MQ error: %v\n", err)
		panic(err)
	}

	channel, e := conn.Channel()
	if e != nil {
		logs.Error("Create channel error: %v\n", e)
		panic(e)
	}

	// 在通道中声明队列
	q, e := channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if e != nil {
		logs.Error("Create queue in channel error: %v\n", e)
		panic(e)
	}

	mq := new(RabbitMQ)
	mq.channel = channel
	mq.Name = q.Name
	return mq
}

// 将队列绑定到交换机
func (q *RabbitMQ) Bind(exchange string) {
	e := q.channel.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,    // no-wait
		nil,      // arguments
	)
	if e != nil {
		panic(e)
	}
	// 将队列绑定到的交换机的名称存储到 q.exchange 中
	q.exchange = exchange
}

// 直接向指定队列中发送消息
func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		logs.Error("Format mq msg error: %v\n", e)
		panic(e)
	}
	e = q.channel.Publish(
		"", // exchange
		queue,
		false,
		false,
		amqp091.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		logs.Error("Send msg to mq error: %v\n", e)
		panic(e)
	}
}

// publish
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		logs.Error("Format mq msg error: %v\n", e)
		panic(e)
	}

	e = q.channel.Publish(
		exchange,
		"", // queue 为空，通过exchange决定发往哪一个queue
		false,
		false,
		amqp091.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		logs.Error("Push msg to mq error: %v\n", e)
		panic(e)
	}

}

// 消费队列中的消息，生成一个接收消息的go channel
func (q *RabbitMQ) Consume() <-chan amqp091.Delivery {
	c, e := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		logs.Error("Consume msg error: %v\n", e)
		panic(e)
	}
	return c
}

// 关闭连接
func (q *RabbitMQ) Close() {
	q.channel.Close()
}
