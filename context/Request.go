package context

import (
	"net/http"
	"net/url"
)

type Request struct {
	request *http.Request
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) Reset(request *http.Request) {
	r.request = request
}

func (r *Request) Request() *http.Request {
	return r.request
}

func (r *Request) RequestUrl() string {
	return r.request.RequestURI
}

func (r *Request) Query() url.Values {
	return r.request.URL.Query()
}

func (r *Request) Method() string {
	return r.request.Method
}
