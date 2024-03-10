package utils

import (
	"data_service/common/config"
	"log"
	"net"
	"strings"
)

// 获取服务当前运行的IP
func GetServerIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// 获取当前HTTP服务的端口
func GetServerHttpPort() string {
	parts := strings.Split(config.Config.Server.Address, ":")
	port := parts[len(parts)-1]
	return port
}
