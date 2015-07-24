package server

import (
	"gopkg.in/go-on/stack.v6"
	"gopkg.in/go-on/stack.v6/mux"
	"gopkg.in/go-on/stack.v6/mw"
	"log"
	"net/http"
)

type Server struct {
	*stack.Stack
	*mux.Mux
}

func New() *Server {
	return &Server{
		Stack: DefaultStack(),
		Mux:   mux.New(),
	}
}

func (s *Server) HTTPServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.Stack.WrapWithContext(stack.HandlerToContextHandler(s.Mux)),
	}
}

func errCatcher(err interface{}, w http.ResponseWriter, r *http.Request) {
	log.Printf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func DefaultStack() *stack.Stack {
	return stack.New().
		Use(mw.Catch(errCatcher)).
		Use(mw.Prepare()).
		Use(mw.MethodOverride()).
		Use(mw.MethodOverrideByField("_method"))
}
