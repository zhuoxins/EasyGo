package context

import (
	"encoding/json"
	"io"
	"net/http"
)

type Output struct {
	ctx         *Context
	SuccessCode int
	ErrorCode   int
	Status      int
}

func NewOutput() *Output {
	return &Output{
		SuccessCode: SUCCESS,
		ErrorCode:   ERROR,
		Status:      http.StatusOK,
	}
}

func (o *Output) Reset(ctx *Context) {
	o.ctx = ctx
}

func (o *Output) Success(args ...interface{}) Render {
	if len(args) == 1 {
		if msg, ok := args[0].(string); ok {
			o.Json(o.SuccessCode, msg, nil)
			return nil
		}
		o.Json(o.SuccessCode, SUCCESSMSG, args[0])
		return nil
	}
	if len(args) == 2 {
		if msg, ok := args[0].(string); ok {
			o.Json(o.SuccessCode, msg, args[1])
			return nil
		}
		o.Json(o.SuccessCode, args[1].(string), args[0])
		return nil
	}
	return &render{code: o.SuccessCode, msg: SUCCESSMSG, out: o}
}

func (o *Output) Error(args ...interface{}) Render {
	if len(args) == 1 {
		if msg, ok := args[0].(string); ok {
			o.Json(o.ErrorCode, msg, nil)
			return nil
		}
		o.Json(o.ErrorCode, ERRORMSG, args[0])
		return nil
	}
	if len(args) == 2 {
		if msg, ok := args[0].(string); ok {
			o.Json(o.ErrorCode, msg, args[1])
			return nil
		}
		o.Json(o.ErrorCode, args[1].(string), args[0])
		return nil
	}
	return &render{code: o.ErrorCode, msg: ERRORMSG, out: o}
}

func (o *Output) Json(code int, msg string, data interface{}, status ...int) {
	result := make(map[string]interface{})
	result = map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	w := o.ctx.Response.Response()
	if len(status) > 0 {
		w.WriteHeader(status[0])
	} else if o.Status != 0 {
		w.WriteHeader(o.Status)
		o.Status = 0
	}
	w.Header().Add(CONTENT_TYPE, CONTENT_JSON)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		o.ctx.Response.Error(err.Error())
	}
}

func (o *Output) WriteString(str string) {
	_, err := io.WriteString(o.ctx.Response.Response(), str)
	if err != nil {
		o.ctx.Response.Error(err.Error())
	}
}

type render struct {
	code int
	msg  string
	data interface{}
	out  *Output
}

func (r *render) WithString(msg string) {
	r.msg = msg
	r.out.Json(r.code, r.msg, nil)
}

func (r *render) WithCode(code int) {
	r.code = code
	r.out.Json(r.code, r.msg, nil)
}

func (r *render) WithData(data interface{}) {
	r.data = data
	r.out.Json(r.code, r.msg, data)
}
