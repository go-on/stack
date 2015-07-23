package mw

import (
	"fmt"
	"net/http"

	"gopkg.in/go-on/stack.v3/responsewriter"
)

type fallback struct {
	handlers    []http.Handler
	ignoreCodes map[int]struct{}
	caller      string
}

func (f *fallback) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {

	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		if _, has := f.ignoreCodes[ck.Code]; !has {
			ck.FlushHeaders()
			ck.FlushCode()
			return true
		}
		ck.Reset()
		return false
	})

	for _, h := range f.handlers {
		h.ServeHTTP(checked, req)
		if checked.HasChanged() {
			if _, has := f.ignoreCodes[checked.Code]; !has {
				checked.FlushMissing()
				return
			}
		}
		checked.Reset()
	}
	next.ServeHTTP(wr, req)
}

func (f *fallback) String() string {
	return fmt.Sprintf("<Fallback %s>", f.caller)
}

// Fallback will try all given handler until
// the first one writes to the ResponseWriter body.
// It is similar to First, but ignores writes that did set one of the ignore status codes
func Fallback(ignoreCodes []int, handler ...http.Handler) *fallback {
	fb := &fallback{handlers: handler, ignoreCodes: map[int]struct{}{}, caller: Caller()}
	for _, code := range ignoreCodes {
		fb.ignoreCodes[code] = struct{}{}
	}
	return fb
}
