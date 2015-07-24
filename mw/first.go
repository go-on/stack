package mw

import (
	"fmt"
	"net/http"

	"gopkg.in/go-on/stack.v4/responsewriter"
)

type first struct {
	handlers []http.Handler
	caller   string
}

func (f first) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})

	for _, h := range f.handlers {
		h.ServeHTTP(checked, req)
		if checked.HasChanged() {
			checked.FlushMissing()
			return
		}
	}
	next.ServeHTTP(wr, req)
}

// First will try all given handler until the first one returns something
func First(handler ...http.Handler) *first {
	return &first{handler, Caller()}
}

func (f *first) String() string {
	return fmt.Sprintf("<First %s>", f.caller)
}
