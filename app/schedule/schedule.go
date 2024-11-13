package main

import (
	"flag"
	"fmt"
	"github.com/xuexihuang/go-SaasSchedule/pkg/data"

	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/config"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/handler"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/schedule.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	err := data.InitMysql(c.Mysql.DataSource, c.Mode)
	if err != nil {
		fmt.Println("InitMysql error!:", err)
		return
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
