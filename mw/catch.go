package mw

import (
	"fmt"
	"net/http"

	"gopkg.in/go-on/stack.v6/responsewriter"
)

type catch struct {
	catchFn func(recovered interface{}, w http.ResponseWriter, r *http.Request)
	caller  string
}

// ServeHTTPNext serves the given request by letting the next serve a ResponseBuffer and
// catching any panics. If no panic happened, the ResponseBuffer is flushed to the ResponseWriter
// Otherwise the CatchFunc is called.
func (c *catch) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})

	defer func() {
		if p := recover(); p != nil {
			c.catchFn(p, wr, req)
		} else {
			checked.FlushMissing()
		}
	}()

	next.ServeHTTP(checked, req)
}

func (c *catch) String() string {
	return fmt.Sprintf("<Catch %s>", c.caller)
}

// Catch returns a CatchFunc for a Catcher
func Catch(fn func(recovered interface{}, w http.ResponseWriter, r *http.Request)) *catch {
	return &catch{fn, Caller()}
}

// PanicCodes returns a middleware that panics if any of the given http codes is set.
// This allows to get stack traces for debugging. Don't use this in production.
func PanicCodes(codes ...int) panicCodes {
	return panicCodes(codes)
}

// panicCodes is a debugging tool to get a stack trace, if some http.StatusCode is set that is inside the list
type panicCodes []int

// ServeHTTP wraps the current Responsewriter with a responsewriter.PanicCodes
func (p panicCodes) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	pn := responsewriter.NewPanicCodes(wr, []int(p)...)
	next.ServeHTTP(pn, req)
	// if we got this far, we had no panic :-)
	pn.FlushAll()
}
