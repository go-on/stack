package mw

import (
	"net/http"

	"github.com/go-on/stack"
	"github.com/go-on/stack/responsewriter"
)

// Error a type based on error that should be saved by a stack.Contexter (response writer)
type Error struct {
	Err error
}

// Swap implements the stack.Swapper interface.
func (e *Error) Swap(repl interface{}) {
	*e = *(repl.(*Error))
}

// SetError stores the given error inside the given stack.Contexter if err is not nil
func SetError(err error, ctx stack.Contexter) {
	if err == nil {
		return
	}
	e := &Error{err}
	ctx.Set(e)
}

// GetError gets the error out of the given stack.Contexter.
// If there is no error, nil is returned.
func GetError(ctx stack.Contexter) error {
	var e Error
	if ctx.Get(&e) {
		return e.Err
	}
	return nil
}

// ErrorHandler is a function that handles an error, if an error has been set inside the stack.Contexter
// via SetError(). It acts as a middleware that passes to the next http.Handler if there is no error inside
// the Contexter, otherwise it handles the error by calling the function
type ErrorHandler func(error, http.ResponseWriter, *http.Request)

func (fn ErrorHandler) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	// returns true, if error happened and was handled, otherwise false
	var handleError = func(rw http.ResponseWriter, req *http.Request) bool {
		err := GetError(rw.(stack.Contexter))
		if err != nil {
			fn(err, rw, req)
			return true
		}
		return false
	}
	bodywritten := false
	checked := responsewriter.NewPeek(rw, func(ck *responsewriter.Peek) bool {
		bodywritten = true
		if handleError(rw, req) {
			return false
		}
		ck.FlushMissing()
		return true
	})

	next.ServeHTTP(checked, req)

	if !bodywritten && !handleError(rw, req) {
		checked.FlushMissing()
	}
}
