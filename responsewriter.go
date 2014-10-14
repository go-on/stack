package stack

import (
	"bufio"
	"net"
	"net/http"
)

// ResponseWriter is used to store the underlying http.ResponseWriter inside the Contexter
type ResponseWriter struct {
	http.ResponseWriter
}

func (rw *ResponseWriter) Swap(repl interface{}) {
	*rw = *(repl.(*ResponseWriter))
}

// ReclaimResponseWriter is a helper that expects the given http.ResponseWriter to either be
// the original http.ResponseWriter or a Contexter which has saved a wrack.ResponseWriter.
// In either case it returns the underlying http.ResponseWriter
func ReclaimResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
	ctx, ok := rw.(Contexter)
	if !ok {
		return rw
	}
	var w ResponseWriter
	if ok := ctx.Get(&w); ok {
		return w.ResponseWriter
	}
	panic("no ResponseWriter saved inside Contexter")
}

// Flush is a helper that flushes the buffer in the  underlying response writer if it is a http.Flusher.
// The http.ResponseWriter might also be a Contexter if it allows the retrieval of the underlying
// ResponseWriter. Ok returns if the underlying ResponseWriter was a http.Flusher
func Flush(rw http.ResponseWriter) (ok bool) {
	w := ReclaimResponseWriter(rw)
	if fl, is := w.(http.Flusher); is {
		fl.Flush()
		return true
	}
	return false
}

// CloseNotify is the same for http.CloseNotifier as Flush is for http.Flusher
// ok tells if it was a CloseNotifier
func CloseNotify(rw http.ResponseWriter) (ch <-chan bool, ok bool) {
	w := ReclaimResponseWriter(rw)
	if cl, is := w.(http.CloseNotifier); is {
		ch = cl.CloseNotify()
		ok = true
		return
	}
	return
}

// Hijack is the same for http.Hijacker as Flush is for http.Flusher
// ok tells if it was a Hijacker
func Hijack(rw http.ResponseWriter) (c net.Conn, brw *bufio.ReadWriter, err error, ok bool) {
	w := ReclaimResponseWriter(rw)
	if hj, is := w.(http.Hijacker); is {
		c, brw, err = hj.Hijack()
		ok = true
		return
	}
	return
}
