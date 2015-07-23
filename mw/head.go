package mw

import (
	"net/http"

	"gopkg.in/go-on/stack.v3/responsewriter"
)

type head struct{}

var _head = head{}

// Head will check if the request method is HEAD.
// if so, it will change the method to GET, call the handler
// with a ResponseBuffer and return only the header information
// to the client.
//
// For non HEAD methods, it simply pass the request handling
// through to the http.handler
func Head() head {
	return _head
}

// ServeHTTPNext handles serves the request by transforming a HEAD request to a GET request
// for the next handler and then remove the body from the response.
// Non HEAD requests are not affected
func (h head) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	if req.Method == "HEAD" {
		req.Method = "GET"

		checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
			ck.FlushHeaders()
			ck.FlushCode()
			return false // write no body
		})

		defer func() {
			req.Method = "HEAD"
		}()

		next.ServeHTTP(checked, req)

		checked.FlushMissing()
		return
	}
	next.ServeHTTP(wr, req)
}
