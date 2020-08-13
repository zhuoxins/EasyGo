package context

import "net/http"

const (
	CONTENT_TYPE    = "Content-Type"
	CONTENT_LENGTH  = "Content-Length"
	CONTENT_BINARY  = "application/octet-stream"
	CONTENT_JSON    = "application/json; charset=utf-8"
	CONTENT_HTML    = "text/html"
	CONTENT_PLAIN   = "text/plain"
	CONTENT_XHTML   = "application/xhtml+xml"
	CONTENT_XML     = "text/xml"
	DEFAULT_CHARSET = "UTF-8"
)

type Response struct {
	ResponseWriter http.ResponseWriter
	Type           string
	Status         int
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Reset(w http.ResponseWriter) {
	r.ResponseWriter = w
	r.Type = "json"
}

func (r *Response) Response() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *Response) Error(errMsg string, code ...int) {
	status := http.StatusInternalServerError
	if len(code) > 0 {
		status = code[0]
	}
	http.Error(r.Response(), errMsg, status)
}

func (r *Response) Header() http.Header {
	header := r.Response().Header()
	return header
}

func (r *Response) SetHeader(k, v string) {
	r.Response().Header().Add(k, v)
}
