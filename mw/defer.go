package mw

import (
	"fmt"
	"net/http"
)

type defer_ struct {
	handler http.Handler
	caller  string
}

func (d *defer_) ServeHTTP(wr http.ResponseWriter, r *http.Request, next http.Handler) {
	defer func() { d.handler.ServeHTTP(wr, r) }()
	next.ServeHTTP(wr, r)
}

func Defer(h http.Handler) *defer_ {
	return &defer_{h, Caller()}
}

func (d *defer_) String() string {
	return fmt.Sprintf("<Defer %s>", d.caller)
}
