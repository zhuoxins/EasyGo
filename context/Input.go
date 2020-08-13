package context

type Input struct {
	ctx         *Context
	requestData map[string]Values
}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Reset(ctx *Context) {
	i.ctx = ctx
	i.requestData = make(map[string]Values)
	var data map[string][]string
	switch i.ctx.Request.Method() {
	case GETREQUEST:
		data = i.ctx.Request.Query()
	case POSTREQUEST:
		_ = i.ctx.Request.Request().ParseMultipartForm(0)
		data = i.ctx.Request.Request().PostForm
	}
	if len(data) > 0 {
		for field, value := range data {
			i.requestData[field] = value
		}
	}
}

func (i *Input) Has(args ...string) bool {
	for _, val := range args {
		if _, ok := i.requestData[val]; !ok {
			return false
		}
	}
	return true
}

//获取单个参数
func (i *Input) Input(field string) Values {
	if val, ok := i.requestData[field]; ok {
		return val
	}
	return nil
}

//获取多个参数
func (i *Input) Only(field ...string) map[string]Values {
	result := make(map[string]Values)
	for _, fields := range field {
		if val, ok := i.requestData[fields]; ok {
			result[fields] = val
		}
	}
	return result
}

func (i *Input) Except(fields ...string) map[string]Values {
	result := make(map[string]Values)
	for field, val := range i.requestData {
		if !inArray(field, fields...) {
			result[field] = val
		}
	}
	return result
}

func (i *Input) All() map[string]Values {
	result := make(map[string]Values)
	for field, val := range i.requestData {
		result[field] = val
	}
	return result
}

func inArray(val string, arr ...string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
