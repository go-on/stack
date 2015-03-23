package stack

import (
	"fmt"
	// "sync"
	// "path/filepath"
	// "runtime"
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

func (s *Stack) ContextHandler() *contextHandler {
	return &contextHandler{s.h}
}

func (s *Stack) String() string {
	var mw []string

	for i, m := range s.Middleware {
		mw = append(mw, fmt.Sprintf("  %p[%d] %s", s, i, m))
	}

	return fmt.Sprintf("<Stack\n  %p %s:%d \n%s\n>", s, s.File, s.Line, strings.Join(mw, "\n"))
}

/*
func NewWithContext(middlewares ...interface{}) *Stack {
	mw := append([]interface{}{Context}, middlewares...)
	s := New(mw...)
	_, file, line, _ := runtime.Caller(1)

	s.Line = line
	s.File = filepath.FromSlash(file)
	return s
}
*/

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
