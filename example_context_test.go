// +build go1.1

package stack_test

import (
	"fmt"
	"net/http"

	"github.com/go-on/stack"
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

func ExampleContext() {
	s := stack.New(
		stack.Context,
		writeVal,
		setVal,
		writeVal,
	)

	r, _ := http.NewRequest("GET", "/", nil)
	s.ServeHTTP(nil, r)

	// Output:
	// no value found
	// value: some value
}
