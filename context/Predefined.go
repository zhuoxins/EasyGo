package context

const (
	GETREQUEST     = "GET"
	POSTREQUEST    = "POST"
	DELETEREQUEST  = "DELETE"
	PATCHREQUEST   = "PATCH"
	OPTIONSREQUEST = "OPTIONS"
	HEADREQUEST    = "HEAD"
	TRACEREQUEST   = "TRACE"
	SUCCESSMSG     = "OK"
	ERRORMSG       = "Error"
)

type Values []string
type Value string

type Render interface {
	WithString(msg string)
	WithCode(code int)
	WithData(data interface{})
}
