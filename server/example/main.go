package main

import (
	"github.com/go-on/stack"
	"github.com/go-on/stack/mw"
	"github.com/go-on/stack/server"
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("did not find anything"))
}

func main() {
	var myStack = stack.New().
		Use(mw.Before(mw.TextString("before\n"))).
		Use(mw.After(mw.TextString("\nafter")))

	fullStack := server.DefaultStack.Concat(myStack)

	s := server.Custom(fullStack)

	s.Router.SetNotFound(http.HandlerFunc(notFound))

	s.INDEX("hu", mw.TextString("hu"))
	s.INDEX("", mw.TextString("index"))

	s.ListenAndServe(":9090")
}
