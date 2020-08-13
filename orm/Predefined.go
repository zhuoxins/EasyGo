package orm

import (
	"EasyGo/config"
	"EasyGo/orm/MysqlOrm"
)

//orm驱动
type driver interface {
	Insert(interface{}) (int, error)
	Delete(...interface{}) (int, error)
	Update(interface{}, ...interface{}) (int, error)
	Fetch(...interface{}) (interface{}, error)
	Query(interface{}, ...interface{}) (interface{}, error)
	ResetConnection(conf config.Result)
}

type Orm interface {
	Insert(data interface{}) (int, error)
	Update(data interface{}, condition ...interface{}) (int, error)
	Delete(condition ...interface{}) (int, error)
	Fetch(condition ...interface{}) (interface{}, error)
	Query(interface{}, ...interface{}) (interface{}, error)
	QueryMS() *MysqlOrm.MysqlOrm
	connect(string, config.Result)
}
