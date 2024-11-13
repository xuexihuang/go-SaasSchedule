package biz

import (
	"fmt"
	log15 "github.com/xuexihuang/new_log15"
)

type AdminNode struct {
	*NodeBase
}

func (a *AdminNode) generateInitSql() string {

	return `CREATE TABLE admin (id INT PRIMARY KEY,name VARCHAR(100));`
}
func (a *AdminNode) generateSetCommand(domain string, imageTag string, tenantId string) ([]string, error) {

	ret := []string{"--set"}
	var setStr string
	sql := a.generateInitSql()
	setStr = "sqlConfig.sql=" + sql
	mysqlUrl := fmt.Sprintf("root123456@tcp(mysql.kube-public.svc.cluster.local3306)%s", tenantId)
	setStr = setStr + ",image.tag=" + imageTag + ",config.Mysql.Database=" + mysqlUrl
	ret = append(ret, setStr)
	log15.Info("generateSetCommand", "ret", ret)
	return ret, nil
}
