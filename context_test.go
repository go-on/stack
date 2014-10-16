package stack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ctx string

func (c *ctx) Swap(i interface{}) {
	*c = *(i.(*ctx))
}

type setCtx string

func (s setCtx) ServeHTTP(c Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	ct := ctx(string(s))
	c.Set(&ct)
	next.ServeHTTP(w, r)
}

func writeCtxNext(c Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	var ct ctx
	c.Get(&ct)
	fmt.Fprint(w, string(ct))
	next.ServeHTTP(w, r)
}

func writeCtx(c Contexter, w http.ResponseWriter, r *http.Request) {
	var ct ctx
	c.Get(&ct)
	fmt.Fprint(w, string(ct))
}

type writeCtx2 struct{}

func (wr writeCtx2) ServeHTTP(c Contexter, w http.ResponseWriter, r *http.Request) {
	var ct ctx
	c.Get(&ct)
	fmt.Fprint(w, string(ct))
}

func delCtx(c Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	var ct ctx
	c.Del(&ct)
	next.ServeHTTP(w, r)
}

type appendCtx string

func (a appendCtx) ServeHTTP(c Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	c.Transaction(func(tc TransactionContexter) {
		var ct ctx
		tc.Get(&ct)
		var n = ctx(string(ct) + string(a))
		tc.Del(&ct) // just to try
		tc.Set(&n)
	})
	next.ServeHTTP(w, r)
}

func TestContext(t *testing.T) {
	s := New(
		Context,
		appendCtx("prepended"),
		setCtx("hiho"),
		Context, // does not do anything
		appendCtx("-appended"),
		writeCtxNext,
		delCtx,
		writeCtx,
		writeCtx2{},
	)

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, nil)

	expected := "hiho-appended"
	if got := rec.Body.String(); got != expected {
		t.Errorf("rec.Body.String() == %#v != %#v", got, expected)
	}
}
