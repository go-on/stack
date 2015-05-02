package mux_test

import (
	"github.com/go-on/stack/mux"
	"github.com/go-on/stack/mw"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux(t *testing.T) {
	h := mux.New()

	tests := []struct {
		mountpoint string
		url        string
		code       int
	}{
		{"/a/", "/a/", http.StatusOK},
		{"/b/", "/b", http.StatusNotFound},
		{"/c/", "/c/x", http.StatusOK},
		{"/d/e/", "/d/e", http.StatusNotFound},
		{"/d/g/", "/f/g/", http.StatusOK},
	}

	for _, test := range tests {
		h.Handle(test.mountpoint, mw.TextString("text"))
	}

	for _, test := range tests {

		rec := httptest.NewRecorder()
		rq, _ := http.NewRequest("GET", test.url, nil)

		h.ServeHTTP(rec, rq)

		if got, want := rec.Code, test.code; got != want {
			t.Errorf("GET %s ; rec.Code = %v; want %v", test.url, got, want)
		}
	}

}
