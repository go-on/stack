package mw

import (
	"fmt"
	"gopkg.in/go-on/stack.v1/responsewriter"
	"net/http"
)

type around struct {
	before, after http.Handler
	caller        string
}

func (a *around) String() string {
	return fmt.Sprintf("<Around %T %v %T %v %s>", a.before, a.before, a.after, a.after, a.caller)
}

func (a around) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	codeOk := false
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		if ck.Code == 200 || ck.Code == 0 {
			codeOk = true
			ck.FlushHeaders()
			a.before.ServeHTTP(wr, req)
			return true
		}
		ck.FlushMissing()
		return true
	})

	next.ServeHTTP(checked, req)
	if codeOk {
		a.after.ServeHTTP(wr, req)
	}
}

// Around returns a wrapper that calls the given before and after handler
// before and after the next handler when serving
func Around(before, after http.Handler) *around {
	return &around{before, after, Caller()}
}
