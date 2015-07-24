package mw

import (
	"net/http"

	"gopkg.in/go-on/stack.v4/responsewriter"
)

type escapeHTML struct{}

func (e escapeHTML) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {
	next.ServeHTTP(responsewriter.NewEscapeHTML(wr), req)
}

// EscapeHTML wraps the next handler by replacing the response writer with an EscapeHTMLResponseWriter
// that escapes html special chars while writing to the underlying response writer
var EscapeHTML = escapeHTML{}
