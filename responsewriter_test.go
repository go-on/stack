// +build go1.1

package stack

import (
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestResponseWriter(t *testing.T) {

	fn := func(ctx Contexter, rw http.ResponseWriter, r *http.Request) {
		w := ReclaimResponseWriter(rw)
		if _, ok := w.(*httptest.ResponseRecorder); ok {
			rw.Write([]byte("ok"))
			return
		}
		rw.Write([]byte("failed"))
	}

	rec := httptest.NewRecorder()
	var s Stack
	s.UseHandlerFuncWithContext(fn)
	s.HandlerWithContext().ServeHTTP(rec, nil)

	expected := "ok"
	if got := rec.Body.String(); got != expected {
		t.Errorf("%#v != %#v", got, expected)
	}
}
