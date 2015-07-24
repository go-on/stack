package mw

import (
	"fmt"
	"gopkg.in/go-on/stack.v4/responsewriter"
	"net/http"
)

type after struct {
	http.Handler
	caller string
}

func After(h http.Handler) *after {
	return &after{h, Caller()}
}

func (b *after) String() string {
	return fmt.Sprintf("<After %T %v %s>", b.Handler, b.Handler, b.caller)
}

func (b *after) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	codeOk := false
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		if ck.Code == 200 || ck.Code == 0 {
			codeOk = true
		}
		ck.FlushMissing()
		return true
	})

	next.ServeHTTP(checked, req)
	if codeOk {
		b.Handler.ServeHTTP(wr, req)
	}
}
