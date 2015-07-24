package main

import (
	"gopkg.in/go-on/stack.v4/mw"
	"gopkg.in/go-on/stack.v4/rest"
	"gopkg.in/go-on/stack.v4/server"
)

func main() {

	s := server.New()
	s.Use(mw.Before(mw.TextString("before\n")))
	s.Use(mw.After(mw.TextString("\nafter")))

	s.Handle("/hu/", rest.Index(mw.TextString("hu")))
	s.Handle("/", rest.Index(mw.TextString("index")))

	s.HTTPServer(":9090").ListenAndServe()
}
