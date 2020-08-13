package config

import (
	"EasyGo/tools/helper"
	"os"
)

//获取配置文件目录
func path(dirName string) string {
	path := helper.Pwd() + helper.Sep() + dirName + helper.Sep()
	if !helper.PathExists(path) {
		_ = os.Mkdir(path, 0755)
	}
	return path
}

type Result interface {
	String(defaultVal ...string) string
	Int() int
	Int64() int64
	Map() map[string]Result
	GetField(field string, defaultVal ...interface{}) Result
	Item(field string, defaultVal ...string) string
}
