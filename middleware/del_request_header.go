package middleware

import (
	"net/http"
	"strings"
)

// DelRequestHeader removes request headers that are identical to the string
// or have it as prefix
type DelRequestHeader string

// ServeHTTPNext removes request headers that are identical to the string
// or have it as prefix. Then the inner http.Handler is called
func (rh DelRequestHeader) ServeHTTPNext(w http.ResponseWriter, r *http.Request, next http.Handler) {
	comp := strings.TrimSpace(strings.ToLower(string(rh)))
	for k := range r.Header {
		k = strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(k, comp) {
			r.Header.Del(k)
		}
	}
	next.ServeHTTP(w, r)
}
