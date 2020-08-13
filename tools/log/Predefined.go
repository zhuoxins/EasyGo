package log

import (
	"EasyGo/tools/helper"
	"os"
)

func path() string {
	path := helper.Pwd() + helper.Sep() + "runtime" + helper.Sep() + "log" + helper.Sep() + helper.Date("Y-m-d")
	if !helper.PathExists(path) {
		_ = os.MkdirAll(path, 0755)
	}
	return path
}
