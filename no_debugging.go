// +build !debug

package stack

import (
	"fmt"
	"net/http"
)

func New(middlewares ...interface{}) *Stack {
	var h http.Handler = http.HandlerFunc(noOp)
	s := &Stack{}

	for i := len(middlewares) - 1; i >= 0; i-- {
		var fn func(http.Handler) http.Handler
		switch x := middlewares[i].(type) {
		case *contextHandler:
			panic("can't embed contextHandler")
		case func(http.Handler) http.Handler:
			fn = x
		case interface {
			Wrap(http.Handler) http.Handler
		}:
			fn = x.Wrap
		case http.Handler:
			fn = end(x.ServeHTTP).middleware
		case func(wr http.ResponseWriter, req *http.Request):
			fn = end(x).middleware
		case func(wr http.ResponseWriter, req *http.Request, next http.Handler):
			fn = mwHandlerFunc(x).Middleware()
		case interface {
			ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler)
		}:
			fn = mwHandler(x)
		case func(ctx Contexter, wr http.ResponseWriter, req *http.Request):
			fn = end(ctxHandlerFunc(x).ServeHTTP).middleware
		case ctxHandler:
			fn = end(ctxHandlerFunc(x.ServeHTTP).ServeHTTP).middleware
		case func(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler):
			fn = mwCtxHandlerFunc(x).Middleware()
		case interface {
			ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
		}:
			fn = mwCtxHandlerFunc(x.ServeHTTP).Middleware()
		default:
			panic(fmt.Sprintf("unsupported middleware type %T, value: %#v", x, x))
		}

		h = fn(h)
	}
	s.h = h
	return s
}
