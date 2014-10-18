package middleware

import (
	"fmt"

	"github.com/go-on/method"
	"github.com/go-on/stack/responsewriter"

	"net/http"
)

type filterBody struct {
	method method.Method
	caller string
}

func (f *filterBody) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if !f.method.Is(r.Method) {
		next.ServeHTTP(w, r)
		return
	}

	checked := responsewriter.NewPeek(w, func(ck *responsewriter.Peek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return false
	})
	next.ServeHTTP(checked, r)

	checked.FlushMissing()
}

// Filter the body for the given method
func FilterBody(m method.Method) *filterBody {
	return &filterBody{m, Caller()}
}

func (f *filterBody) String() string {
	return fmt.Sprintf("<FilterBody %s>", f.caller)
}
