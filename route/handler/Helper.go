package handler

import (
	"EasyGo/route"
	"strings"
)

type helper struct {
	router *routeHandler
}

func NewHelper(router *routeHandler) route.RouterHelper {
	return &helper{router: router}
}

func (h *helper) Filter(filters []string) route.RouterHelper {
	for _, v := range filters {
		if v != "" {
			parseRegex := strings.Split(v, "=>")
			h.router.Progress.Regex[strings.TrimSpace(parseRegex[0])] = strings.TrimSpace(parseRegex[1])
		}
	}
	return h
}
func (h *helper) PreMiddleware(exec route.Middleware) route.RouterHelper {
	h.router.Progress.AddPreMiddleware(exec)
	return h
}

func (h *helper) BackMiddleware(exec route.Middleware) route.RouterHelper {
	h.router.Progress.AddBackMiddleware(exec)
	return h
}
