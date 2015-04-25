// +build go1.1

package stack

import (
	"net/http"
)

func New() (s *Stack) {
	return &Stack{}
}

// Stack is a stack of middlewares that handle http requests
type Stack []func(http.Handler) http.Handler

// Use adds the given middleware to the middleware stack
func (s *Stack) Use(mw Middleware) *Stack {
	return s.UseWrapperFunc(mwHandler(mw))
}

// UseFunc adds the given function to the middleware stack
func (s *Stack) UseFunc(fn func(wr http.ResponseWriter, req *http.Request, next http.Handler)) *Stack {
	return s.UseWrapperFunc(mwHandlerFunc(fn).Middleware())
}

// UseHandler adds the given handler as middleware to the stack.
// the handler will be called before the next middleware
func (s *Stack) UseHandler(mw http.Handler) *Stack {
	return s.UseWrapper(before(mw.ServeHTTP))
}

// UseHandlerFunc is like UseHandler but for http.HandlerFunc
func (s *Stack) UseHandlerFunc(fn func(wr http.ResponseWriter, req *http.Request)) *Stack {
	return s.UseWrapper(before(fn))
}

// UseHandlerWithContext adds the given context handler as middleware to the stack.
// the handler will be called before the next middleware
func (s *Stack) UseHandlerWithContext(mw ContextHandler) *Stack {
	return s.UseWrapper(before(ContextHandlerFunc(mw.ServeHTTP).ServeHTTP))
}

// UseHandlerFuncWithContext adds the given function as middleware to the stack.
// the handler will be called before the next middleware
func (s *Stack) UseHandlerFuncWithContext(fn func(c Contexter, w http.ResponseWriter, r *http.Request)) *Stack {
	return s.UseWrapper(before(ContextHandlerFunc(fn).ServeHTTP))
}

// UseWithContext adds the context middleware to the middleware stack
func (s *Stack) UseWithContext(mw ContextMiddleware) *Stack {
	return s.UseWrapperFunc(mwCtxHandlerFunc(mw.ServeHTTP).Middleware())
}

// UseFuncWithContext adds the given function to the middleware stack
func (s *Stack) UseFuncWithContext(fn func(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)) *Stack {
	return s.UseWrapperFunc(mwCtxHandlerFunc(fn).Middleware())
}

// UseWrapper adds the given wrapper to the middleware stack
func (s *Stack) UseWrapper(mw Wrapper) *Stack {
	return s.UseWrapperFunc(mw.Wrap)
}

// UseWrapperFunc adds the given function to the middleware stack
func (s *Stack) UseWrapperFunc(mw func(http.Handler) http.Handler) *Stack {
	*s = append(*s, mw)
	return s
}

// Concat returns a new stack that has the middleware of the current stack concatenated with the middleware of the given stack
func (s *Stack) Concat(st *Stack) *Stack {
	// *s = append(*s, st...)
	s3 := append(*s, (*st)...)
	return &s3
}

// Wrap wraps the stack around the next handler and returns the resulting handler
func (s *Stack) Wrap(next http.Handler) http.Handler {
	return s.wrap(next)
}

func (s *Stack) wrap(next http.Handler) http.Handler {
	if next == nil {
		next = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}
	for i := len(*s) - 1; i >= 0; i-- {
		next = (*s)[i](next)
	}
	return next
}

func (s *Stack) Handler() http.Handler {
	return s.wrap(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
}

func (s *Stack) HandlerWithContext() http.Handler {
	return &contextHandler{s.wrap(nil)}
}

func (s *Stack) WrapFunc(fn func(wr http.ResponseWriter, req *http.Request)) http.Handler {
	return s.Wrap(http.HandlerFunc(fn))
}

// WrapWithContext wraps the stack around the given app and returns a handler to run the stack
// that adds context to the http.ResponseWriter.
// It should be used instead of Wrap for the outermost stack and only there
func (s *Stack) WrapWithContext(app ContextHandler) http.Handler {
	return &contextHandler{s.Wrap(ContextHandlerFunc(app.ServeHTTP))}
}

func (s *Stack) WrapFuncWithContext(fn func(ctx Contexter, wr http.ResponseWriter, req *http.Request)) http.Handler {
	return &contextHandler{s.Wrap(ContextHandlerFunc(fn))}
}
