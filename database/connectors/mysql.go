package connectors

import (
	"EasyGo/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//连接mysql数据库
func MysqlConnect(configs config.Result) (*sql.DB, error) {
	dsn := configs.Item("username") + ":" + configs.Item("password") + "@tcp(" + configs.Item("host") + ":" + configs.Item("port") + ")/" + configs.Item("database") + "?"
	//设置字符集
	if charset := configs.Item("charset"); charset != "" {
		dsn += "&charset=" + charset
	} else {
		dsn += "&charset=utf8"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		if db != nil {
			db.Close()
		}
		return nil, err
	}
	//设置最大连接
	if maxConns := configs.GetField("maxConns").Int(); maxConns != 0 {
		db.SetMaxIdleConns(maxConns)
	}
	//最大闲时连接
	if maxOpenConns := configs.GetField("maxOpenConns").Int(); maxOpenConns != 0 {
		db.SetMaxOpenConns(maxOpenConns)
	}
	return db, err
}
