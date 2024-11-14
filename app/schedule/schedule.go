package main

import (
	"flag"
	"fmt"
	"github.com/xuexihuang/go-SaasSchedule/pkg/data"
	"log"
	"os"

	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/config"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/handler"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gopkg.in/yaml.v2"
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
	fmt.Println("initmysql success==", c.Mysql.DataSource)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

type Paths struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}
type Hosts struct {
	Host  string   `yaml:"host"`
	Paths []*Paths `yaml:"paths"`
}

// 定义结构体
type UserConfig struct {
	Image struct {
		Tag string `yaml:"tag"`
	}
	Ingress struct {
		Enabled bool     `yaml:"enabled"`
		Hosts   []*Hosts `yaml:"hosts"`
	}
	Config struct {
		Host  string `yaml:"Host"`
		Port  int    `yaml:"Port"`
		Mysql struct {
			DataSource string `yaml:"Port"`
		}
	}
	SqlConfig struct {
		Sql string `yaml:"sql"`
	}
}

func main2() {
	c := UserConfig{}
	c.Ingress.Enabled = true
	c.Ingress.Hosts = make([]*Hosts, 1)
	c.Ingress.Hosts[0] = &Hosts{Host: "www.game36666.com", Paths: []*Paths{&Paths{Path: "/user(/|$)(.*)", PathType: "ImplementationSpecific"}}}
	c.Image.Tag = "v1.0"
	c.Config.Host = "0.0.0.0:80"
	c.Config.Port = 80
	c.Config.Mysql.DataSource = `root:123456@tcp(172.30.33.164:30306)/saas_schedule?charset=utf8mb4&parseTime=true`
	c.SqlConfig.Sql = `-- Example SQL content
CREATE TABLE user (
id INT PRIMARY KEY,
name VARCHAR(100)
);
insert user(0,"jack");
insert user(0,"huanglin");
insert user(0,"tom");`

	// 序列化为 YAML
	data, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 打印 YAML 输出
	fmt.Println(string(data))

	// 将 YAML 数据写入到文件
	err = os.WriteFile("config.yaml", data, 0644)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}

	fmt.Println("YAML 文件已成功生成")
}
