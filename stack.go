package stack

import (
	"fmt"
	"runtime"
	"strings"

	"net/http"
)

var noOp = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

type Middleware interface {
	ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler)
}

type ContextHandler interface {
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request)
}

type ContextMiddleware interface {
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
}

type Wrapper interface {
	Wrap(http.Handler) http.Handler
}

type Swapper interface {
	// Swap must be defined on a pointer
	// and changes the the value of the pointer
	// to the value the replacement is pointing to
	Swap(replacement interface{})
}

/*
	Accepted middlewares

	Functions:

	func(wr http.ResponseWriter, req *http.Request)
	func(ctx Contexter, wr http.ResponseWriter, req *http.Request)
	func(wr http.ResponseWriter, req *http.Request, next http.Handler)
	func(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
	func(http.Handler) http.Handler

	Interfaces:

	http.Handler: ServeHTTP(wr http.ResponseWriter, req *http.Request)
	ContextHandler: ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request)
	Middleware: ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler)
	ContextMiddleware: ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
	Wrapper: Wrap(http.Handler) http.Handler
*/

type Stack struct {
	h          http.Handler
	Middleware []string
	File       string
	Line       int
}

func (s *Stack) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.h.ServeHTTP(w, req)
}

func (s *Stack) String() string {
	var mw []string

	for i, m := range s.Middleware {
		mw = append(mw, fmt.Sprintf("  %p[%d] %s", s, i, m))
	}

	return fmt.Sprintf("<Stack\n  %p %s:%d \n%s\n>", s, s.File, s.Line, strings.Join(mw, "\n"))
}

func New(middlewares ...interface{}) *Stack {
	var h http.Handler = http.HandlerFunc(noOp)
	s := &Stack{}
	s.Middleware = make([]string, len(middlewares))
	_, s.File, s.Line, _ = runtime.Caller(1)

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

type mwHandlerFunc func(wr http.ResponseWriter, req *http.Request, next http.Handler)

func (m mwHandlerFunc) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) { m(wr, req, next) })
	}
}

func mwHandler(h interface {
	ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler)
}) func(http.Handler) http.Handler {
	return mwHandlerFunc(h.ServeHTTP).Middleware()
}

type ctxHandler interface {
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request)
}

type ctxHandlerFunc func(ctx Contexter, wr http.ResponseWriter, req *http.Request)

func (c ctxHandlerFunc) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	ctx, ok := wr.(Contexter)
	if !ok {
		panic("wrack.Context not in stack")
	}
	c(ctx, wr, req)
}

type mwCtxHandlerFunc func(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)

func (m mwCtxHandlerFunc) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			ctx := wr.(Contexter)
			m(ctx, wr, req, next)
		})
	}
}

func mwCtxHandler(h interface {
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
}) func(http.Handler) http.Handler {
	return mwCtxHandlerFunc(h.ServeHTTP).Middleware()
}

/*
	Converts to an http.Handler

	The following functions are converted:

	func(http.Handler) http.Handler
	func(wr http.ResponseWriter, req *http.Request)
	func(ctx Contexter, wr http.ResponseWriter, req *http.Request)
	func(wr http.ResponseWriter, req *http.Request, next http.Handler)
	func(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)

	The following interfaces are converted:

	Wrap(http.Handler) http.Handler
	ServeHTTP(wr http.ResponseWriter, req *http.Request)
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request)
	ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler)
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
*/
func Handler(i interface{}) (h http.Handler) {
	switch x := i.(type) {
	case http.Handler:
		h = x
	case http.HandlerFunc:
		h = x
	case func(wr http.ResponseWriter, req *http.Request):
		h = http.HandlerFunc(x)
	case ctxHandlerFunc:
		h = x
	case func(ctx Contexter, wr http.ResponseWriter, req *http.Request):
		h = ctxHandlerFunc(x)
	case ctxHandler:
		h = ctxHandlerFunc(x.ServeHTTP)
	case func(wr http.ResponseWriter, req *http.Request, next http.Handler):
		h = mwHandlerFunc(x).Middleware()(noOp)
	case func(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler):
		h = mwCtxHandlerFunc(x).Middleware()(noOp)
	case interface {
		ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler)
	}:
		h = mwHandlerFunc(x.ServeHTTP).Middleware()(noOp)
	case interface {
		ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
	}:
		h = mwCtxHandlerFunc(x.ServeHTTP).Middleware()(noOp)
	case func(http.Handler) http.Handler:
		h = x(noOp)
	default:
		panic(fmt.Sprintf("can't convert %T to handler", x))
	}

	return
}

func HandlerFunc(i interface{}) (h http.HandlerFunc) {
	return Handler(i).ServeHTTP
}

type end http.HandlerFunc

func (e end) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	http.HandlerFunc(e).ServeHTTP(wr, req)
}

func (e end) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(e)
}

func Handlers(i ...interface{}) (h []http.Handler) {
	h = make([]http.Handler, len(i))

	for n, ii := range i {
		h[n] = Handler(ii)
	}

	return h
}
