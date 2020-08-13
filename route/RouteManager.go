package route

import (
	"EasyGo/context"
	"net/http"
	"reflect"
	"sync"
)

type RouteManager struct {
	Router Router
	pool   sync.Pool
}

//自定义http服务管理
func NewManager(router Router) *RouteManager {
	manage := &RouteManager{
		Router: router,
		pool: sync.Pool{
			New: func() interface{} {
				return context.NewContext()
			},
		},
	}
	return manage
}

//接收请求
func (s *RouteManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := s.pool.Get().(*context.Context)
	ctx.Reset(w, r)
	//转发doHandler
	s.doHandler(ctx)
}

//匹配请求
func (s *RouteManager) doHandler(ctx *context.Context) {
	defer s.pool.Put(ctx)
	info, err := s.Router.MatchRoute(ctx)
	if err != nil {
		ctx.Response.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	info = s.Router.Extract(info, ctx)
	ctx.Route.Reset(ctx, info.RouteParam)
	//执行前置中间件
	if len(info.PreMiddleware) > 0 {
		for _, middleware := range info.PreMiddleware {
			middleware(ctx)
		}
	}
	//执行相关方法
	if info.HandleFunc != "" {
		//执行实例
		execInstance := reflect.New(reflect.TypeOf(info.Instance).Elem())
		//判断是否有构造方法
		initFunc := execInstance.MethodByName("Init")
		if initFunc.String() != "<invalid Value>" {
			initFunc.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
		//判断控制器中是否有该方法
		execMethod := execInstance.MethodByName(info.HandleFunc)
		if execMethod.String() == "<invalid Value>" {
			//请求方法不存在
			ctx.Response.Error("route: request method not exists", http.StatusInternalServerError)
			return
		}
		execMethod.Call(make([]reflect.Value, 0))
	} else {
		//直接执行方法体路由
		if execMethod, ok := info.Instance.(FuncRequest); ok {
			execMethod(ctx)
		}
	}
	//执行后置中间件
	if len(info.BackMiddleware) > 0 {
		for _, middleware := range info.BackMiddleware {
			middleware(ctx)
		}
	}
}
