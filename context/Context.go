package context

import "net/http"

type Context struct {
	Input    *Input
	Output   *Output
	Request  *Request
	Response *Response
	Route    *Route
}

func NewContext() *Context {
	return &Context{
		Input:    NewInput(),
		Output:   NewOutput(),
		Request:  NewRequest(),
		Response: NewResponse(),
		Route:    NewRoute(),
	}
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.Request.Reset(r)
	c.Response.Reset(w)
	c.Input.Reset(c)
	c.Output.Reset(c)
}
