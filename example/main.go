package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-on/stack"
	"github.com/go-on/stack/mw"
)

type ctxA string

func (c *ctxA) Swap(v interface{}) {
	*c = *(v.(*ctxA))
}

func setctxA(ctx stack.Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler) {
	c := ctxA(fmt.Sprintf("hi ctx %p", req))
	ctx.Set(&c)
	next.ServeHTTP(wr, req)
}

func printctxA(ctx stack.Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler) {
	var c ctxA
	ctx.Get(&c)

	fmt.Fprintf(wr, "(ctxA: %#v vs %p)", c, req)
}

type s string

func (ss s) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	fmt.Fprint(wr, string(ss))
}

func hello(wr http.ResponseWriter, req *http.Request) {
	fmt.Fprint(wr, "hello")
}

func main() {

	w := stack.New(
		mw.Before(s("hi")),
		setctxA,
		mw.Before(
			stack.New(
				mw.Before(s("inner1")),
				mw.Before(s("inner2")),
				printctxA,
			),
		),
		mw.Before(stack.Handler(hello)),
		s("end"),
	)

	h := w.ContextHandler()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	fmt.Println(w)
}
