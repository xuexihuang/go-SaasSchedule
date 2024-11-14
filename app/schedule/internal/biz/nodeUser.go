package biz

import (
	"fmt"
)

type UserNode struct {
	*NodeBase
}

func (u *UserNode) generateInitSql() string {
	return `-- Example SQL content
CREATE TABLE user (
id INT PRIMARY KEY,
name VARCHAR(100)
);
insert user(0,"jack");
insert user(0,"huanglin");
insert user(0,"tom");`
}
func (u *UserNode) generateSetCommand(domain string, imageTag string, tenantId string) (interface{}, error) {

	c := UserConfig{}
	c.Ingress.Enabled = true
	c.Ingress.Hosts = make([]*Hosts, 1)
	c.Ingress.Hosts[0] = &Hosts{Host: domain, Paths: []*Paths{&Paths{Path: "/userapi(/|$)(.*)", PathType: "ImplementationSpecific"}}}
	c.Image.Tag = imageTag
	c.Config.Host = "0.0.0.0:80"
	c.Config.Port = 80
	c.Config.Mysql.DataSource = fmt.Sprintf("root:123456@tcp(172.30.33.164:30306)/%s?charset=utf8mb4&parseTime=true", tenantId)
	c.SqlConfig.Sql = u.generateInitSql()
	return &c, nil

}
