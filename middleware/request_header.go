package middleware

import "net/http"

type requestHeader struct {
	Key, Val string
}

// SetRequestHeader sets a request header
func RequestHeader(key, val string) *requestHeader {
	return &requestHeader{key, val}
}

func (s *requestHeader) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	req.Header.Set(s.Key, s.Val)
	next.ServeHTTP(wr, req)
}
