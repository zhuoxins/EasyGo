package context

type Route struct {
	ctx  *Context
	data map[string]Value
}

func NewRoute() *Route {
	return &Route{}
}

func (this *Route) Reset(ctx *Context, data map[string]string) {
	this.ctx = ctx
	this.data = make(map[string]Value)
	for k, v := range data {
		this.data[k] = Value(v)
	}
}

//获取单个参数
func (this *Route) Input(field string) Value {
	if val, ok := this.data[field]; ok {
		return val
	}
	return ""
}

func (this *Route) Get(fields ...string) map[string]Value {
	result := make(map[string]Value)
	for _, field := range fields {
		if val, ok := this.data[field]; ok {
			result[field] = val
		}
	}
	return result
}

func (this *Route) All() map[string]Value {
	result := make(map[string]Value)
	for key, val := range this.data {
		result[key] = val
	}
	return result
}

func (v Value) Strings() string {
	return string(v)
}
