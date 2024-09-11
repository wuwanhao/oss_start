package heartbeat

import (
	"common_service/logs"
	"connector_service/rabbitmq"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// 缓存所有的数据服务节点
var dataServers = make(map[string]time.Time)

// 使用互斥锁保护dataServers这个map的并发读写
var mutex sync.Mutex

// 接收和处理来自数据服务节点的心跳消息
func ListenHeartbeat() {
	q := rabbitmq.New()
	defer q.Close()

	q.Bind("apiServers")
	c := q.Consume()

	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

// 10S后剔除
func removeExpiredDataServer() {
	for {
		time.Sleep(time.Second)
		mutex.Lock()
		for dataServer, lastLiveTime := range dataServers {
			if lastLiveTime.Add(10 * time.Second).Before(time.Now()) {
				logs.Info("remove data node: %v\n", dataServer)
				delete(dataServers, dataServer)
			}
		}
		mutex.Unlock()
	}
}

// 获取可用的数据服务列表
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// 随机选择一个数据服务
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
