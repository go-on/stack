package middleware

import (
	"fmt"
	"net/http"
	"runtime"
)

func Caller() (info string) {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}

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
	//http.HandlerFunc(b).ServeHTTP(wr, req)
	b.Handler.ServeHTTP(wr, req)
	next.ServeHTTP(wr, req)
}
