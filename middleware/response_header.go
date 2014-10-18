package middleware

import "net/http"

type responseHeader struct {
	Key, Val string
}

// ResponseHeader sets a response header
func ResponseHeader(key, val string) *responseHeader {
	return &responseHeader{key, val}
}

// ServeHTTPNext sets the header of the key to the value and calls
// the next handler after that
func (s *responseHeader) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	wr.Header().Set(s.Key, s.Val)
	next.ServeHTTP(wr, req)
}
