package handler

import "EasyGo/route"

type groupRouting struct {
	router *routeHandler
}

func BuildGroup(router *routeHandler) *groupRouting {
	return &groupRouting{router: router}
}

//路由分组
func (this *groupRouting) Group(instance interface{}, routeSet route.RouteSeter) {
	defer this.router.ReGroup()
	this.router.IsGroup = true
	this.router.GroupInstance = instance
	routeSet(this.router)
}

func (this *groupRouting) Prefix(prefix string) route.GroupRouter {
	this.router.Prefix = prefix
	return this
}

func (this *groupRouting) PreMiddleware(exec route.Middleware) route.GroupRouter {
	this.router.AddPreMiddleware(exec)
	return this
}

func (this *groupRouting) BackMiddleware(exec route.Middleware) route.GroupRouter {
	this.router.AddBackMiddleware(exec)
	return this
}
