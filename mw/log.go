package mw

import (
	"log"
	"net/http"
)

// a simple request logger
type logger struct{}

func (l logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	log.Printf("%s %s", r.Method, r.URL.String())
	next.ServeHTTP(w, r)
}

var Logger = logger{}
