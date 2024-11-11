package data

import (
	_ "github.com/go-sql-driver/mysql"
	log15 "github.com/xuexihuang/new_log15"
	"github.com/zeromicro/go-zero/core/service"
	"xorm.io/xorm"
)

var MysqlEngine *xorm.Engine

func InitMysql(dataSource, mode string) error {
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log15.Error("init mysql Engine error", "err", err)
		return err
	}
	if mode == service.DevMode {
		engine.ShowSQL(true)
	}
	MysqlEngine = engine
	return nil
}

// 在事务中执行的方法handler
type transHandler func(session *xorm.Session) (interface{}, error)

// DoInTrans 在事务中执行指定方法：f
func DoInTrans(f transHandler) (interface{}, error) {
	return MysqlEngine.Transaction(f)
}

func NewSession() *xorm.Session {
	return MysqlEngine.NewSession()
}
