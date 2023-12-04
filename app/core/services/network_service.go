package services

import (
	"context"
	"fmt"
	"net"
	"scrm-sync-data/app/config"
	"strings"
	"time"
)

type NetworkService struct {
	context context.Context
	cancel  context.CancelFunc
	config  *config.Config //日志配置
	State   bool           //服务状态
}

func (networkService *NetworkService) Init(parentContext context.Context, config *config.Config) {
	networkService.config = config
	networkService.State = false
	networkService.context, networkService.cancel = context.WithCancel(parentContext)
	networkService.Start()
}

func (networkService *NetworkService) Start() {
	fmt.Println("start network service...", time.Now())
	networkService.State = true
}

func (networkService *NetworkService) Stop() {
	fmt.Println("stop network service...", time.Now())
	networkService.cancel()
	networkService.State = false
}

func (networkService *NetworkService) GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
