// +build go1.1

package stack_test

import (
	"fmt"
	"net/http"

	"github.com/go-on/stack"
)

/*
This example illustrates 3 ways to write and use middleware.

For sharing context, look at example_context_test.go.
*/

type print1 string

func (p print1) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	fmt.Print(p)
	next.ServeHTTP(wr, req)
}

type print2 string

func (p print2) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		fmt.Print(p)
		next.ServeHTTP(rw, req)
	}
	return f
}

type print3 string

func (p print3) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	fmt.Println(p)
}

func ExampleNew() {
	s := stack.New(
		print1("ready..."),
		print2("steady..."),
		print3("go!"),
		// if there should be a handler after this, you will need the Before wrapper from go-on/wrack/middleware
	)
	r, _ := http.NewRequest("GET", "/", nil)
	s.ServeHTTP(nil, r)

	// Output:
	// ready...steady...go!
	//
}
