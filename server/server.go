package server

import (
	"github.com/go-on/stack"
	"github.com/go-on/stack/mw"
	"github.com/go-on/stack/router"
	"log"
	"net/http"
)

type Server struct {
	stack *stack.Stack
	*router.Router
}

func New() *Server {
	return Custom(DefaultStack)
}

func Custom(st *stack.Stack) *Server {
	return &Server{
		Router: router.New(),
		stack:  st,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.stack.WrapWithContext(nil))
}

func errCatcher(err interface{}, w http.ResponseWriter, r *http.Request) {
	log.Printf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

var DefaultStack = stack.New().
	Use(mw.Catch(errCatcher)).
	Use(mw.Prepare()).
	Use(mw.MethodOverride()).
	Use(mw.MethodOverrideByField("_method"))
