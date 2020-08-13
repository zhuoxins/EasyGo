package handler

import (
	"EasyGo/context"
	"EasyGo/kernel"
	"EasyGo/route"
	helper2 "EasyGo/tools/helper"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

func init() {
	kernel.Register("route").Do(func() {
		router, grouper := NewHandler()
		route.RouteBuilder = router
		route.Manage = route.NewManager(router)
		route.GroupHandler = grouper
	})
}

//请求方式
const (
	GETREQUEST     = "GET"
	POSTREQUEST    = "POST"
	DELETEREQUEST  = "DELETE"
	PUTREQUEST     = "PUT"
	PATCHREQUEST   = "PATCH"
	OPTIONSREQUEST = "OPTIONS"
	ANYREQUEST     = "ANY"
)

type routeHandler struct {
	routeTree            map[string]map[string]*route.RouteInfo
	regexTree            map[string]map[string]*route.RouteInfo
	Progress             *route.RouteInfo
	Prefix               string
	IsGroup              bool
	GroupInstance        interface{}
	globalPreMiddleware  []route.Middleware
	globalBackMiddleware []route.Middleware
}

//返回路由操作对象
func NewHandler() (route.Router, route.GroupRouter) {
	handler := &routeHandler{
		IsGroup:   false,
		routeTree: make(map[string]map[string]*route.RouteInfo),
		regexTree: make(map[string]map[string]*route.RouteInfo),
	}
	groupHandler := BuildGroup(handler)
	return handler, groupHandler
}

func (this *routeHandler) AddPreMiddleware(handle route.Middleware) {
	this.globalPreMiddleware = append(this.globalPreMiddleware, handle)
}

func (this *routeHandler) AddBackMiddleware(handle route.Middleware) {
	this.globalBackMiddleware = append(this.globalBackMiddleware, handle)
}

func (this *routeHandler) reset() {
	this.Progress = &route.RouteInfo{
		Regex:        make(map[string]string),
		DefaultParam: make(map[string]string),
	}
}

//初始化执行实例
func (this *routeHandler) initInstance(execInstance ...interface{}) {
	this.reset()
	if !this.IsGroup && len(execInstance) == 0 {
		panic("route: miss handle instance")
	}
	if this.IsGroup && len(execInstance) == 0 {
		this.Progress.Instance = this.GroupInstance
	}
	if len(execInstance) > 0 {
		this.Progress.Instance = execInstance[0]
	}
}

//执行方法体路由
func (this *routeHandler) Route(path string, instance route.FuncRequest, reqMethod ...string) route.RouterHelper {
	this.initInstance(instance)
	method := ANYREQUEST
	if len(reqMethod) > 0 {
		method = strings.ToUpper(strings.ToLower(reqMethod[0]))
	}
	this.appendMethod(method, path, "")
	return NewHelper(this)
}

//GET request
func (this *routeHandler) Get(path, execMethod string, execInstance ...interface{}) route.RouterHelper {
	this.initInstance(execInstance...)
	this.appendMethod(GETREQUEST, path, execMethod)
	return NewHelper(this)
}

//POST request
func (this *routeHandler) Post(path, execMethod string, execInstance ...interface{}) route.RouterHelper {
	this.initInstance(execInstance...)
	this.appendMethod(POSTREQUEST, path, execMethod)
	return NewHelper(this)
}

//DELETE request
func (this *routeHandler) Delete(path, execMethod string, execInstance ...interface{}) route.RouterHelper {
	this.initInstance(execInstance...)
	this.appendMethod(DELETEREQUEST, path, execMethod)
	return NewHelper(this)
}

//PUT request
func (this *routeHandler) Put(path, execMethod string, execInstance ...interface{}) route.RouterHelper {
	this.initInstance(execInstance...)
	this.appendMethod(PUTREQUEST, path, execMethod)
	return NewHelper(this)
}

//PATCH request
func (this *routeHandler) Patch(path, execMethod string, execInstance ...interface{}) route.RouterHelper {
	this.initInstance(execInstance...)
	this.appendMethod(PATCHREQUEST, path, execMethod)
	return NewHelper(this)
}

//anyway request
func (this *routeHandler) Any(path, execMethod string, execInstance ...interface{}) route.RouterHelper {
	this.initInstance(execInstance...)
	this.appendMethod(ANYREQUEST, path, execMethod)
	return NewHelper(this)
}

////restful 资源路由
//func (this *routeHandler) Resource(path string, execInstance ...interface{}) route.Router {
//	return this
//}

//重置分组信息
func (this *routeHandler) ReGroup() {
	this.IsGroup = false
	this.rePrefix()
	this.reMiddleware()
}

//重置路由前缀
func (this *routeHandler) rePrefix() {
	this.Prefix = ""
}

//重置中间件信息
func (this *routeHandler) reMiddleware() {
	this.globalPreMiddleware = make([]route.Middleware, 0)
	this.globalBackMiddleware = make([]route.Middleware, 0)
}

//注册路由信息
func (this *routeHandler) appendMethod(reqMethod, path, method string) {
	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}
	if this.Prefix != "" {
		prefix := strings.TrimRight(this.Prefix, "/")
		if strings.Index(prefix, "/") != 0 {
			prefix = "/" + prefix
		}
		path = prefix + path
	}
	regex, defaultParam := this.disposeRouteParams(path)
	this.Progress.Path = path
	this.Progress.Method = reqMethod
	if len(defaultParam) > 0 {
		this.Progress.DefaultParam = defaultParam
	}
	this.Progress.HandleFunc = method

	if len(this.globalPreMiddleware) > 0 {
		for _, handle := range this.globalPreMiddleware {
			this.Progress.AddPreMiddleware(handle)
		}
	}
	if len(this.globalBackMiddleware) > 0 {
		for _, handle := range this.globalBackMiddleware {
			this.Progress.AddBackMiddleware(handle)
		}
	}
	if regex == "" {
		this.Progress.Type = route.NORMALROUTING
		if tree, ok := this.routeTree[reqMethod]; ok {
			tree[path] = this.Progress
			this.routeTree[reqMethod] = tree
		} else {
			this.routeTree[reqMethod] = map[string]*route.RouteInfo{
				path: this.Progress,
			}
		}
	} else {
		this.Progress.Type = route.PARAMROUTING
		if tree, ok := this.regexTree[reqMethod]; ok {
			tree[regex] = this.Progress
			this.regexTree[reqMethod] = tree
		} else {
			this.regexTree[reqMethod] = map[string]*route.RouteInfo{
				regex: this.Progress,
			}
		}
	}
}

//判断是否有路由参数
func (this *routeHandler) disposeRouteParams(path string) (string, map[string]string) {
	regex, defaultParam := "", make(map[string]string)
	if ok := strings.Contains(path, ":"); ok {
		parsePath := strings.Split(path, "/")
		for _, v := range parsePath {
			if v == "" {
				continue
			}
			if isParamRoute(v) {
				regex += `\/[\s\S][^/]*?`
			} else if isDefaultRoute(v) {
				regex += `[\s\S]*?`
				parse := strings.Split(helper2.TrimStr(v, "{", "}"), ",")
				if !helper2.InSlice(1, parse) {
					panic("route : missing default value")
				}
				defaultParam[strings.TrimLeft(parse[0], ":")] = strings.TrimSpace(parse[1])
			} else {
				regex += `\/` + v
			}
		}
	}
	return regex, defaultParam
}

//提取路由参数
func (this *routeHandler) Extract(info *route.ControllerInfo, ctx *context.Context) *route.ControllerInfo {
	//提取路由参数
	if info.Type == route.PARAMROUTING {
		parsePath, parseUrl := slicePath(info.Path, ctx.Request.RequestUrl())
		info.RouteParam = make(map[string]string)
		for k, v := range parsePath {
			if isParamRoute(v) {
				info.RouteParam[strings.TrimLeft(v, ":")] = parseUrl[k]
			} else if isDefaultRoute(v) {
				paramName := extractDefName(v)
				if helper2.InSlice(k, parseUrl) {
					info.RouteParam[paramName] = parseUrl[k]
				} else {
					info.RouteParam[paramName] = info.DefaultParam[paramName]
				}
			}
		}
	}
	return info
}

/**
根据请求方式匹配相应路由
*/
func (this *routeHandler) MatchRoute(ctx *context.Context) (*route.ControllerInfo, error) {
	var routeInfo *route.RouteInfo
	path := ctx.Request.RequestUrl()
	if strings.Contains(path, "?") {
		path = strings.Split(path, "?")[0]
	}
	reqMethod := ctx.Request.Method()
	var wg sync.WaitGroup
	wg.Add(2)
	//匹配无参路由
	go func() {
		defer wg.Done()
		if tree, ok := this.routeTree[reqMethod]; ok {
			if info, ok := tree[path]; ok {
				routeInfo = info
				return
			}
		}
		//没有继续匹配Any请求
		if tree, ok := this.routeTree[ANYREQUEST]; ok {
			if info, ok := tree[path]; ok {
				routeInfo = info
				return
			}
		}
	}()
	//匹配有参路由
	go func() {
		defer wg.Done()
		if tree, ok := this.regexTree[reqMethod]; ok {
			if ok, info := matchParamRoute(tree, path); ok {
				routeInfo = info
				return
			}
		}
		//没有继续匹配Any请求
		if tree, ok := this.regexTree[ANYREQUEST]; ok {
			if ok, info := matchParamRoute(tree, path); ok {
				routeInfo = info
				return
			}
		}
	}()
	wg.Wait()
	if routeInfo != nil {
		controllerType := reflect.TypeOf(routeInfo.Instance)
		controller := &route.ControllerInfo{
			Type:           routeInfo.Type,
			Path:           routeInfo.Path,
			Instance:       routeInfo.Instance,
			HandleFunc:     routeInfo.HandleFunc,
			ControllerType: controllerType,
		}
		controller.Initialize(routeInfo)
		return controller, nil
	}
	return nil, errors.New("route not found")
}

func matchParamRoute(tree map[string]*route.RouteInfo, path string) (bool, *route.RouteInfo) {
	for regex, info := range tree {
		if ok, _ := regexp.MatchString(regex, path); ok {
			if len(info.Regex) > 0 {
				parseUrl := strings.Split(path, "/")[1:]
				parsePath := strings.Split(info.Path, "/")[1:]
				for k, v := range parsePath {
					var urlParam, paramRegex string
					if isParamRoute(v) {
						paramRegex, ok = info.Regex[v]
						//判断是否有自定义正则
						if !ok || paramRegex == "" {
							continue
						}
						urlParam = parseUrl[k]
					} else if isDefaultRoute(v) {
						//默认参数值路由
						paramName := strings.Split(helper2.TrimStr(v, "{", "}"), ",")[0]
						paramRegex, ok = info.Regex[paramName]
						if !ok || paramRegex == "" {
							continue
						}
						//如果url有参数,取url参数,否则取默认参数值
						if helper2.InSlice(k, parseUrl) {
							urlParam = parseUrl[k]
						} else {
							urlParam, _ = info.DefaultParam[strings.TrimLeft(paramName, ":")]
						}
					}
					if ok, _ := regexp.MatchString(paramRegex, urlParam); ok {
						continue
					} else {
						return false, nil
					}
				}
			}
			return true, info
		}
	}
	return false, nil
}

//判断是否有参路由
func isParamRoute(path string) bool {
	if strings.Contains(path, ":") && strings.Index(path, ":") == 0 {
		return true
	}
	return false
}

//判断是否有默认值路由
func isDefaultRoute(path string) bool {
	if strings.Index(path, "{") == 0 && strings.Index(path, "}") == (len(path)-1) && strings.Index(path, ":") == 1 {
		return true
	}
	return false
}

//分割路由
func slicePath(path, urlPath string) (parsePath []string, parseUrl []string) {
	parsePath = strings.Split(path, "/")
	parseUrl = strings.Split(urlPath, "/")
	return
}

//提取默认路由参数名
func extractDefName(str string) string {
	name := strings.TrimLeft(strings.Split(helper2.TrimStr(str, "{", "}"), ",")[0], ":")
	return name
}
