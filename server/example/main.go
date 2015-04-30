package main

import (
	"github.com/go-on/stack/mw"
	"github.com/go-on/stack/rest"
	"github.com/go-on/stack/server"
)

func main() {

	s := server.New()
	s.Use(mw.Before(mw.TextString("before\n")))
	s.Use(mw.After(mw.TextString("\nafter")))

	s.Handle("/hu/", rest.Index(mw.TextString("hu")))
	s.Handle("/", rest.Index(mw.TextString("index")))

	s.HTTPServer(":9090").ListenAndServe()
}
