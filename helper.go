package stack

import (
	"net/http"
)

type mwHandlerFunc func(wr http.ResponseWriter, req *http.Request, next http.Handler)

func (m mwHandlerFunc) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) { m(wr, req, next) })
	}
}

func mwHandler(h Middleware) func(http.Handler) http.Handler {
	return mwHandlerFunc(h.ServeHTTP).Middleware()
}

func HandlerToContextHandler(h http.Handler) ContextHandler {
	return ctxHandlerFunc(func(ctx Contexter, wr http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(wr, req)
	})
}

type ctxHandlerFunc func(ctx Contexter, wr http.ResponseWriter, req *http.Request)

func (c ctxHandlerFunc) ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request) {
	c(ctx, wr, req)
}

type ContextHandlerFunc func(ctx Contexter, wr http.ResponseWriter, req *http.Request)

func (c ContextHandlerFunc) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
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

func mwCtxHandler(h ContextMiddleware) func(http.Handler) http.Handler {
	return mwCtxHandlerFunc(h.ServeHTTP).Middleware()
}

type before http.HandlerFunc

func (b before) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	http.HandlerFunc(b).ServeHTTP(wr, req)
	next.ServeHTTP(wr, req)
}

func (b before) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		http.HandlerFunc(b).ServeHTTP(wr, req)
		next.ServeHTTP(wr, req)
	})
}
