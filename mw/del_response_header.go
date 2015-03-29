package mw

import (
	"net/http"
	"strings"

	"gopkg.in/go-on/wrap.v2"
)

// DelResponseHeader removes response headers that are identical to the string
// or have if as prefix
type DelResponseHeader string

// ServeHTTP removes the response headers that are identical to the string
// or have if as prefix after the next handler is run
func (rh DelResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	bodyWritten := false
	comp := strings.TrimSpace(strings.ToLower(string(rh)))

	checked := wrap.NewPeek(w, func(ck *wrap.Peek) bool {
		hd := ck.Header()
		for k := range hd {
			k = strings.TrimSpace(strings.ToLower(k))
			if strings.HasPrefix(k, comp) {
				hd.Del(k)
			}
		}
		ck.FlushHeaders()
		ck.FlushCode()
		bodyWritten = true
		return true
	})

	next.ServeHTTP(checked, r)

	if !bodyWritten {
		hd := checked.Header()
		for k := range hd {
			k = strings.TrimSpace(strings.ToLower(k))
			if strings.HasPrefix(k, comp) {
				hd.Del(k)
			}
		}
	}
	checked.FlushMissing()
}
