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

func ExampleStack() {
	var s stack.Stack
	s.Use(print1("ready..."))
	s.UseWrapper(print2("steady..."))
	s.UseHandler(print3("go!"))
	r, _ := http.NewRequest("GET", "/", nil)
	s.Handler().ServeHTTP(nil, r)

	// Output:
	// ready...steady...go!
	//
}
