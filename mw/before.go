package mw

import (
	"fmt"
	"gopkg.in/go-on/stack.v3/responsewriter"
	"net/http"
)

type before struct {
	http.Handler
	caller string
}

func Before(h http.Handler) *before {
	return &before{h, Caller()}
}

func (b *before) String() string {
	return fmt.Sprintf("<Before %T %v %s>", b.Handler, b.Handler, b.caller)
}

func (b *before) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		if ck.Code == 200 || ck.Code == 0 {
			ck.FlushHeaders()
			b.Handler.ServeHTTP(wr, req)
			return true
		}
		ck.FlushMissing()
		return true
	})

	next.ServeHTTP(checked, req)
}
