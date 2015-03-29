// +build go1.1

package stack_test

import (
	"fmt"
	"net/http"

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

func exampleDebug() {

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

	fmt.Println(w)

	// Output:
	// <Stack
	// 0xc20803e140 /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:55
	// 0xc20803e140[0] func(http.Handler) http.Handler = (func(http.Handler) http.Handler)(0x445e90)
	// 0xc20803e140[1] *mw.before = <Before stack_test.s hi /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:44>
	// 0xc20803e140[2] func(stack.Contexter, http.ResponseWriter, *http.Request, http.Handler) = (func(stack.Contexter, http.ResponseWriter, *http.Request, http.Handler))(0x44ba50)
	// 0xc20803e140[3] *mw.before = <Before *stack.Stack <Stack
	// 0xc20803e100 /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:51
	// 0xc20803e100[0] *mw.before = <Before stack_test.s inner1 /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:48>
	// 0xc20803e100[1] *mw.before = <Before stack_test.s inner2 /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:49>
	// 0xc20803e100[2] func(stack.Contexter, http.ResponseWriter, *http.Request, http.Handler) = (func(stack.Contexter, http.ResponseWriter, *http.Request, http.Handler))(0x44bbf0)
	// > /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:52>
	// 0xc20803e140[4] *mw.before = <Before http.HandlerFunc 0x44bf00 /home/benny/Entwicklung/gopath/src/github.com/go-on/stack/example_debug_test.go:53>
	// 0xc20803e140[5] stack_test.s = "end"
	// >

}
