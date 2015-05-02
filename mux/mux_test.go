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
		text       string
	}{
		{"/a/", "/a/", http.StatusOK, "A"},
		{"/b/", "/b/c", http.StatusOK, "B"},
		{"/c/", "/c/x", http.StatusOK, "C"},
		{"/d/e/", "/d/e", http.StatusNotFound, ""},
		{"/d/g/", "/d/g/", http.StatusOK, "DG"},
		{"/d/p", "/d/p", http.StatusOK, "DP"},
		{"/d/h", "/d/", http.StatusNotFound, ""},
		{"/e/h", "/e/h", http.StatusOK, "EH"},
		{"/e/", "/e/f", http.StatusOK, "E"},
		{"/e/g", "/e/g", http.StatusOK, "EG"},
	}

	for _, test := range tests {
		h.Handle(test.mountpoint, mw.TextString(test.text))
	}

	for _, test := range tests {

		rec := httptest.NewRecorder()
		rq, _ := http.NewRequest("GET", test.url, nil)

		h.ServeHTTP(rec, rq)

		if got, want := rec.Code, test.code; got != want {
			t.Errorf("GET %s ; rec.Code = %v; want %v", test.url, got, want)
		}

		if rec.Code == 200 {
			if got, want := rec.Body.String(), test.text; got != want {
				t.Errorf("GET %s ; rec.Body.String() = %v; want %v", test.url, got, want)
			}
		}
	}

}
