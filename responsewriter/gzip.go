package responsewriter

// modelled after code of Andrew Gerrand as proposed here https://groups.google.com/forum/#!searchin/golang-nuts/http%2420gzip/golang-nuts/eVnTcMwNVjM/u0a6TQLagnkJ

import (
	"compress/gzip"
	"io"
	"net/http"

	"gopkg.in/go-on/stack.v1"
)

type GZip struct {
	io.WriteCloser
	http.ResponseWriter
	stack.Contexter
}

func (w *GZip) Write(b []byte) (int, error) {
	return w.WriteCloser.Write(b)
}

func NewGZIP(rw http.ResponseWriter) *GZip {
	gz := &GZip{WriteCloser: gzip.NewWriter(rw), ResponseWriter: rw}
	if ctx, ok := rw.(stack.Contexter); ok {
		gz.Contexter = ctx
	}
	return gz
}
