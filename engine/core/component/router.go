package core

import (
	"net/http"
)

type Router struct {
	routes     map[string]ComponentInterface
	middleware []MiddlewareFunc
}

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

func NewRouter() *Router {
	return &Router{
		routes:     make(map[string]ComponentInterface),
		middleware: make([]MiddlewareFunc, 0),
	}
}

func (r *Router) AddRoute(path string, component ComponentInterface) {
	r.routes[path] = component
}
