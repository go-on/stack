// +build go1.1

package stack_test

import (
	"fmt"
	"net/http"

	"github.com/go-on/stack"
)

type MyVal string

// Swap replaces the value of m with the value of val
func (m *MyVal) Swap(val interface{}) {
	*m = *(val.(*MyVal)) // will never panic
}

type MyMiddleware struct{}

func (mw *MyMiddleware) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	var val MyVal = "some value"
	ctx.Set(&val)
	next.ServeHTTP(rw, req)
}

func writeVal(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	var val MyVal
	found := ctx.Get(&val)
	if !found {
		fmt.Println("no value found")
	} else {
		fmt.Printf("value: %s\n", string(val))
	}
	next.ServeHTTP(rw, req)
}

func ExampleContext() {
	rack := stack.New(
		stack.Context,
		writeVal,
		&MyMiddleware{},
		writeVal,
	)

	r, _ := http.NewRequest("GET", "/", nil)
	rack.ServeHTTP(nil, r)

	// Output:
	// no value found
	// value: some value
}
