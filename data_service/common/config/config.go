package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// 总配置
type config struct {
	Server   server   `yaml:"server"`
	RabbitMq rabbitmq `yaml:"rabbitmq"`
	Oss      oss      `yaml:"oss"`
}

// 服务配置
type server struct {
	Address string `yaml:"address"`
	Name    string `yaml:"name"`
	Mode    string `yaml:"mode"`
}

// OSS配置
type oss struct {
	StorageRoot  string `yaml:"storage-root"`
	StorageIndex string `yaml:"storage-index"`
}

// mq配置
type rabbitmq struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Virtualhost  string `yaml:"vhost"`
	ApiExchange  string `yaml:"api-exchange"`
	DataExchange string `yaml:"data-exchange"`
}

// 用于对外暴露的配置对象
var Config *config

// 配置初始化
func init() {
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
