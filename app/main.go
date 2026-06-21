package main

import (
	"github.com/setcreed/hade-kit/config"
	"github.com/setcreed/hade-kit/logs"
	"github.com/setcreed/hade-kit/server"

	"app/internal/inits"
)

func main() {
	//1. 加载配置  默认是 etc/config.yml
	config.Init()
	conf := config.GetConfig()
	//2. 加载日志
	logs.Init(conf.Log)
	s := server.NewServer(conf)
	//4. 初始化模块
	inits.Init(s, conf)
	s.Start()
}
