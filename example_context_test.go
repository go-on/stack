// +build go1.1

package stack_test

import (
	"fmt"
	"github.com/go-on/stack"
	"net/http"
)

type Val string

// Swap replaces the value of m with the value of val
func (m *Val) Swap(val interface{}) {
	*m = *(val.(*Val)) // will never panic
}

func setVal(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	var val Val = "some value"
	ctx.Set(&val)
	next.ServeHTTP(rw, req)
}

func writeVal(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	var val Val
	found := ctx.Get(&val)
	if !found {
		fmt.Println("no value found")
	} else {
		fmt.Printf("value: %s\n", string(val))
	}
	next.ServeHTTP(rw, req)
}

func ExampleContexter() {
	var s stack.Stack
	s.UseFuncWithContext(writeVal)
	s.UseFuncWithContext(setVal)
	s.UseFuncWithContext(writeVal)

	r, _ := http.NewRequest("GET", "/", nil)
	s.HandlerWithContext().ServeHTTP(nil, r)

	// Output:
	// no value found
	// value: some value
}
