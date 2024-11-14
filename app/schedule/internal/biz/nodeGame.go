package biz

import "fmt"

type GameNode struct {
	*NodeBase
}

func (u *GameNode) generateInitSql() string {
	return `-- Example SQL content
CREATE TABLE game (
id INT PRIMARY KEY,
name VARCHAR(100)
);
insert game(0,"jack");
insert game(0,"huanglin");
insert game(0,"tom");`
}
func (u *GameNode) generateSetCommand(domain string, imageTag string, tenantId string) (interface{}, error) {

	c := UserConfig{}
	c.Ingress.Enabled = true
	c.Ingress.Hosts = make([]*Hosts, 1)
	c.Ingress.Hosts[0] = &Hosts{Host: domain, Paths: []*Paths{&Paths{Path: "/gameapi(/|$)(.*)", PathType: "ImplementationSpecific"}}}
	c.Image.Tag = imageTag
	c.Config.Host = "0.0.0.0:80"
	c.Config.Port = 80
	c.Config.Mysql.DataSource = fmt.Sprintf("root:123456@tcp(172.30.33.164:30306)/%s?charset=utf8mb4&parseTime=true", tenantId)
	c.SqlConfig.Sql = u.generateInitSql()
	return &c, nil

}
