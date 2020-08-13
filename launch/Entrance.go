package launch

import (
	_ "EasyGo/cache"
	_ "EasyGo/config"
	"EasyGo/kernel"
	"EasyGo/route"
	_ "EasyGo/tools/builder"
)

func Run() {
	//加载服务
	kernel.Load()
	//路由监听
	route.Listen()
}
