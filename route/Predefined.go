package route

import (
	"EasyGo/config"
	"EasyGo/context"
	"log"
	"net/http"
	"reflect"
)

//路由类型
const (
	NORMALROUTING = iota
	PARAMROUTING
)

var (
	//http请求管理者
	Manage *RouteManager
	//路由服务
	RouteBuilder Router
	GroupHandler GroupRouter
)

//中间件
type Middleware func(ctx *context.Context)
type FuncRequest func(ctx *context.Context)
type RouteSeter func(r Router)

//路由组
type GroupRouter interface {
	Group(instance interface{}, routeSet RouteSeter)
	PreMiddleware(exec Middleware) GroupRouter
	BackMiddleware(exec Middleware) GroupRouter
	Prefix(prefix string) GroupRouter
}

//定义路由管理
type Router interface {
	Route(path string, instance FuncRequest, reqMethod ...string) RouterHelper
	Get(path, execMethod string, execInstance ...interface{}) RouterHelper
	Post(path, execMethod string, execInstance ...interface{}) RouterHelper
	Delete(path, execMethod string, execInstance ...interface{}) RouterHelper
	Put(path, execMethod string, execInstance ...interface{}) RouterHelper
	Patch(path, execMethod string, execInstance ...interface{}) RouterHelper
	Any(path, execMethod string, execInstance ...interface{}) RouterHelper
	MatchRoute(ctx *context.Context) (*ControllerInfo, error)
	Extract(info *ControllerInfo, ctx *context.Context) *ControllerInfo
}

//路由助手
type RouterHelper interface {
	Filter(filters []string) RouterHelper
	PreMiddleware(exec Middleware) RouterHelper
	BackMiddleware(exec Middleware) RouterHelper
}

//路由控制器
type RouteController interface {
	Init()
}

//路由信息
type RouteInfo struct {
	Type           int
	Path           string
	Method         string
	Regex          map[string]string
	DefaultParam   map[string]string
	Instance       interface{}
	HandleFunc     string
	PreMiddleware  []Middleware
	BackMiddleware []Middleware
}

func (r *RouteInfo) AddPreMiddleware(handle Middleware) {
	r.PreMiddleware = append(r.PreMiddleware, handle)
}

func (r *RouteInfo) AddBackMiddleware(handle Middleware) {
	r.BackMiddleware = append(r.BackMiddleware, handle)
}

//路由请求
type ControllerInfo struct {
	Type           int
	Path           string
	Instance       interface{}
	HandleFunc     string
	ControllerType reflect.Type
	RouteParam     map[string]string
	DefaultParam   map[string]string
	PreMiddleware  []Middleware
	BackMiddleware []Middleware
}

func (this *ControllerInfo) Initialize(route *RouteInfo) {
	if len(route.DefaultParam) > 0 {
		this.DefaultParam = make(map[string]string)
		for k, v := range route.DefaultParam {
			this.DefaultParam[k] = v
		}
	}
	if len(route.PreMiddleware) > 0 {
		for _, middleware := range route.PreMiddleware {
			this.PreMiddleware = append(this.PreMiddleware, middleware)
		}
	}
	if len(route.BackMiddleware) > 0 {
		for _, middleware := range route.BackMiddleware {
			this.BackMiddleware = append(this.BackMiddleware, middleware)
		}
	}
}

//监听http端口
func Listen() {
	portConf := config.Get("app.basis.port").String()
	port := ":"
	if portConf == "" {
		portConf = "8421"
	}
	port += portConf
	err := http.ListenAndServe(port, Manage)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

/**
 *路由快捷函数
 */
func Group(instance interface{}, routeSet RouteSeter) {
	GroupHandler.Group(instance, routeSet)
}

func Get(path, execMethod string, execInstance ...interface{}) RouterHelper {
	return RouteBuilder.Get(path, execMethod, execInstance...)
}

func Post(path, execMethod string, execInstance ...interface{}) RouterHelper {
	return RouteBuilder.Post(path, execMethod, execInstance...)
}

func Delete(path, execMethod string, execInstance ...interface{}) RouterHelper {
	return RouteBuilder.Delete(path, execMethod, execInstance...)
}

func Put(path, execMethod string, execInstance ...interface{}) RouterHelper {
	return RouteBuilder.Put(path, execMethod, execInstance...)
}

func Patch(path, execMethod string, execInstance ...interface{}) RouterHelper {
	return RouteBuilder.Patch(path, execMethod, execInstance...)
}

func Any(path, execMethod string, execInstance ...interface{}) RouterHelper {
	return RouteBuilder.Any(path, execMethod, execInstance...)
}

func Prefix(prefix string) GroupRouter {
	return GroupHandler.Prefix(prefix)
}

func PreMiddleware(m Middleware) GroupRouter {
	return GroupHandler.PreMiddleware(m)
}

func BackMiddleware(m Middleware) GroupRouter {
	return GroupHandler.BackMiddleware(m)
}

func Route(path string, f FuncRequest, reqMethod ...string) RouterHelper {
	return RouteBuilder.Route(path, f, reqMethod...)
}
