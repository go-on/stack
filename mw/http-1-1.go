package mw

import "net/http"

type http1_1 struct{}

var HTTP1_1 = http1_1{}

func (h http1_1) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if !r.ProtoAtLeast(1, 1) {
		// protocol not supported
		w.WriteHeader(505)
		return
	}
	next.ServeHTTP(w, r)
}
