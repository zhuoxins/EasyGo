package database

import (
	"EasyGo/config"
	"EasyGo/database/connectors"
)

func Connect(driver string, conf config.Result) (conn interface{}, err error) {
	switch driver {
	case "mysql":
		return connectors.MysqlConnect(conf)
	default:
		conn, err = nil, nil
	}
	return
}
