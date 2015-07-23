package mw

import (
	"fmt"
	"net/http"

	"gopkg.in/go-on/stack.v1/responsewriter"
)

// GuardFunc is a responsewriter.Wapper and http.HandlerFunc that may operate on the ResponseWriter
// If it does so, the responsewriterper prevents the next Handler from serving.
type guard struct {
	http.Handler
	caller string
}

// ServeHTTPNext lets the GuardFunc serve to a ResponseBuffer and if it changed something
// the Response is send to the ResponseWriter, preventing the next http.Handler from
// executing. Otherwise the next handler serves the request.
func (g *guard) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})
	g.Handler.ServeHTTP(checked, req)
	if checked.HasChanged() {
		checked.FlushMissing()
		return
	}
	next.ServeHTTP(wr, req)
}

func (g *guard) String() string {
	return fmt.Sprintf("<Guard %T %v %s>", g.Handler, g.Handler, g.caller)
}

// Guard returns a GuardFunc for a http.Handler
func Guard(h http.Handler) *guard {
	return &guard{h, Caller()}
}
