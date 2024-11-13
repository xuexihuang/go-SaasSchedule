package biz

import "fmt"

type UserNode struct {
	*NodeBase
}

func (u *UserNode) generateInitSql() string {
	return `CREATE TABLE admin (id INT PRIMARY KEY,name VARCHAR(100));`
}
func (u *UserNode) generateSetCommand(domain string, imageTag string, tenantId string) ([]string, error) {

	ret := []string{"--set"}
	var setStr string
	sql := u.generateInitSql()
	setStr = "sqlConfig.userSql=" + sql
	mysqlUrl := fmt.Sprintf("root:123456@tcp(mysql.kube-public.svc.cluster.local:3306)/%s?charset=utf8mb4&parseTime=true", tenantId)
	setStr = setStr + ",image.tag=" + imageTag + ",config.Mysql.Database=" + mysqlUrl
	ret = append(ret, setStr)
	return ret, nil

}
