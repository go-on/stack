package mw

import (
	"net/http"

	"github.com/go-on/stack/responsewriter"
)

// ContentType writes the content type if the next handler was successful
// and did not set a content-type
type ContentType string

// SetContentType sets the content type in the given ResponseWriter
func (c ContentType) Set(w http.ResponseWriter) {
	w.Header().Set("Content-Type", string(c))
}

// ServeHandle serves the given request with the next handler and after that
// writes the content type, if the next handler was successful
// and did not set a content-type
func (c ContentType) ServeHTTP(wr http.ResponseWriter, req *http.Request, next http.Handler) {

	var bodyWritten = false
	checked := responsewriter.NewPeek(wr, func(ck *responsewriter.Peek) bool {
		if ck.IsOk() {
			c.Set(ck)
		}
		ck.FlushHeaders()
		ck.FlushCode()
		bodyWritten = true
		return true
	})

	next.ServeHTTP(checked, req)

	if !bodyWritten {
		c.Set(checked)
	}
	checked.FlushMissing()
}

var (
	JSONContentType       = ContentType("application/json; charset=utf-8")
	TextContentType       = ContentType("text/plain; charset=utf-8")
	CSSContentType        = ContentType("text/css; charset=utf-8")
	HTMLContentType       = ContentType("text/html; charset=utf-8")
	JavaScriptContentType = ContentType("application/javascript; charset=utf-8")
	RSSFeedContentType    = ContentType("application/rss+xml; charset=utf-8")
	AtomFeedContentType   = ContentType("application/atom+xml; charset=utf-8")
	PDFContentType        = ContentType("application/pdf")
	CSVContentType        = ContentType("text/csv; charset=utf-8")
)
