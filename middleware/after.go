package middleware

import (
	"fmt"
	"net/http"
)

type after struct {
	http.Handler
	caller string
}

func After(h http.Handler) *after {
	return &after{h, Caller()}
}

func (b *after) String() string {
	return fmt.Sprintf("<After %T %v %s>", b.Handler, b.Handler, b.caller)
}

func (b *after) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	next.ServeHTTP(wr, req)
	b.Handler.ServeHTTP(wr, req)
}
