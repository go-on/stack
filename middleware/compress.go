package middleware

import (
	"net/http"
	"strings"

	"github.com/go-on/stack/responsewriter"
)

// GZip compresses the body written by the next handlers
// on the fly if the client did set the Accept-Encoding header to gzip.
// It also sets the response header Content-Encoding to gzip if it did the compression.
func GZip(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		next.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := responsewriter.NewGZIP(w)
	defer gz.Close()
	next.ServeHTTP(gz, r)
}

// Deflate compresses the body written by the next handlers with the level of the int
type Deflate int

func (d Deflate) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		next.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "deflate")
	df := responsewriter.NewDeflate(w, int(d))
	defer df.Close()
	next.ServeHTTP(df, r)
}

// Compress compresses the body written by the next handlers via either gzip or deflate, if the
// client supports it. When deflate is used the int is used as compression level
type Compress int

func (d Compress) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := responsewriter.NewGZIP(w)
		defer gz.Close()
		next.ServeHTTP(gz, r)
		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		w.Header().Set("Content-Encoding", "deflate")
		df := responsewriter.NewDeflate(w, int(d))
		defer df.Close()
		next.ServeHTTP(df, r)
		return
	}

	next.ServeHTTP(w, r)
}
