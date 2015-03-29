package mw

import (
	"net/http"
)

// casts the Responsewriter to http.Handler in order to write to itself, not sure if it is of use
type responseWriterHandle struct{}

func (rwh responseWriterHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.(http.Handler).ServeHTTP(w, r)
}

var ResponseWriterHandler = responseWriterHandle{}
