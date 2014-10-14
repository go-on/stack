package responsewriter

// modelled after code of Andrew Gerrand as proposed here https://groups.google.com/forum/#!searchin/golang-nuts/http%2420gzip/golang-nuts/eVnTcMwNVjM/u0a6TQLagnkJ

import (
	"compress/flate"
	"io"
	"net/http"

	"github.com/go-on/stack"
)

type Deflate struct {
	io.WriteCloser
	http.ResponseWriter
	stack.Contexter
}

func (w *Deflate) Write(b []byte) (int, error) {
	return w.WriteCloser.Write(b)
}

// if level < 0, default compression is used
// if level = 0, no compression is used
// if level >= 9, max compression is used
func NewDeflate(rw http.ResponseWriter, level int) *Deflate {

	// no error for level [-1, 9]
	if level < -1 {
		level = -1
	}

	if level > 9 {
		level = 9
	}

	fl, _ := flate.NewWriter(rw, level)
	df := &Deflate{WriteCloser: fl, ResponseWriter: rw}
	if ctx, ok := rw.(stack.Contexter); ok {
		df.Contexter = ctx
	}
	return df
}
