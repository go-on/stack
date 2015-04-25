// +build go1.1

package stack

import (
	"net/http"
	"testing"
)

func TestStack(t *testing.T) {

	tests := map[string]http.Handler{
		"abc": New().
			UseHandler(write("a")).
			UseFunc(write("b").ServeHTTPNext).
			UseHandlerFunc(write("c").ServeHTTP).
			Handler(),
		"ab": New().
			UseWrapperFunc(write("a").Wrap).
			UseWrapper(writeStop("b")).
			UseHandler(write("c")).
			Handler(),
	}

	for body, h := range tests {
		rec, req := newTestRequest("GET", "/")
		h.ServeHTTP(rec, req)
		assertResponse(t, rec, body, 200)
	}
}
