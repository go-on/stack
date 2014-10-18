package middleware

import (
	"fmt"
	"net/http"
)

type before struct {
	http.Handler
	caller string
}

func Before(h http.Handler) *before {
	return &before{h, Caller()}
}

func (b *before) String() string {
	return fmt.Sprintf("<Before %T %v %s>", b.Handler, b.Handler, b.caller)
}

func (b *before) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	b.Handler.ServeHTTP(wr, req)
	next.ServeHTTP(wr, req)
}
