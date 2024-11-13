package biz

import (
	"fmt"
	log15 "github.com/xuexihuang/new_log15"
)

type UserNode struct {
	*NodeBase
}

func (u *UserNode) generateInitSql() string {
	return `CREATE TABLE admin`
}
func (u *UserNode) generateSetCommand(domain string, imageTag string, tenantId string) ([]string, error) {

	ret := []string{"--set"}
	var setStr string
	sql := u.generateInitSql()
	setStr = "sqlConfig.sql=" + sql
	mysqlUrl := fmt.Sprintf("root123456tcpmysql.kube-public.svc.cluster.local3306-%s", tenantId)
	setStr = setStr + ",image.tag=" + imageTag + ",config.Mysql.Database=" + mysqlUrl
	ret = append(ret, setStr)
	log15.Info("generateSetCommand", "ret", ret)
	return ret, nil

}
