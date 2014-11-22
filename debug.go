// +build debug

package stack

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

func init() {
	// fmt.Println("running init")
	New = func(middlewares ...interface{}) *Stack {
		var h http.Handler = http.HandlerFunc(noOp)
		s := &Stack{}
		s.Middleware = make([]string, len(middlewares))
		_, file, line, _ := runtime.Caller(1)

		s.Line = line
		s.File = filepath.FromSlash(file)

		for i := len(middlewares) - 1; i >= 0; i-- {
			var debugInfo string
			if insp, ok := middlewares[i].(fmt.Stringer); ok {
				debugInfo = fmt.Sprintf("%T = %s", middlewares[i], insp.String())
			} else {
				debugInfo = fmt.Sprintf("%T = %#v", middlewares[i], middlewares[i])
			}
			s.Middleware[i] = debugInfo
			var fn func(http.Handler) http.Handler
			switch x := middlewares[i].(type) {
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

}
