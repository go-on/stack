package mw

import (
	"fmt"
	"net/http"

	"github.com/go-on/stack/responsewriter"
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
