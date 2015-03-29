package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mustPanic(t *testing.T, fn func()) {
	defer func() {
		e := recover()
		if e == nil {
			t.Errorf("must panic")
		}
	}()
	fn()
}

func TestMountRessourceTwice(t *testing.T) {
	rt := New()
	rs := &Ressource{}
	rs.Mount("a", rt)

	mustPanic(t, func() {
		rs.Mount("a", rt)
	})

}

func TestRessource(t *testing.T) {

	tests := []struct {
		handlerType   string
		requestMethod string
		requestURL    string
		pathSegment   string
		mountPoint    string
		called        bool
		code          int
		id            string
		url           string
	}{
		{"GET", "GET", "/a/1", "a", "", true, 200, "1", "/a/1"},
		{"POST", "POST", "/b", "b", "", true, 200, "", "/b"},
		{"PUT", "PUT", "/c/1", "c", "", true, 200, "1", "/c/1"},
		{"PATCH", "PATCH", "/d/1", "d", "", true, 200, "1", "/d/1"},
		{"DELETE", "DELETE", "/e/1", "e", "", true, 200, "1", "/e/1"},
		{"INDEX", "GET", "/f/", "f", "", true, 200, "", "/f"},
	}

	for _, test := range tests {
		rt := New()
		rs := &Ressource{}
		fh := &fakeHandler{}
		hf := fh.ServeHTTP

		switch test.handlerType {
		case "GET":
			rs.GET = hf
		case "POST":
			rs.POST = hf
		case "PATCH":
			rs.PATCH = hf
		case "PUT":
			rs.PUT = hf
		case "DELETE":
			rs.DELETE = hf
		case "INDEX":
			rs.INDEX = hf
		}

		rs.Mount(test.pathSegment, rt)

		var url string
		switch test.handlerType {
		case "GET":
			url = rs.GetURL(test.id)
		case "POST":
			url = rs.PostURL()
		case "PATCH":
			url = rs.PatchURL(test.id)
		case "PUT":
			url = rs.PutURL(test.id)
		case "DELETE":
			url = rs.DeleteURL(test.id)
		case "INDEX":
			url = rs.IndexURL()
		}

		rq, err := http.NewRequest(test.requestMethod, test.requestURL, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		rec := httptest.NewRecorder()

		if test.mountPoint != "" {
			parent := New()
			rt.Mount(test.mountPoint, parent)
			parent.ServeHTTP(rec, rq)
		} else {
			rt.ServeHTTP(rec, rq)
		}

		identifier := test.requestMethod + " " + test.requestURL + " for " + test.pathSegment

		if test.mountPoint != "" {
			identifier += " mounted at " + test.mountPoint
		}

		if got, want := url, test.url; got != want {
			t.Errorf("%s url = %v; want %v", identifier, got, want)
		}

		if got, want := rec.Code, test.code; got != want {
			t.Errorf("%s rec.Code = %v; want %v", identifier, got, want)
		}

		if got, want := fh.called, test.called; got != want {
			t.Errorf("%s fh.called = %v; want %v", identifier, got, want)
		}

		if test.id != "" && fh.called {
			if got, want := rq.URL.Fragment, test.id; got != want {
				t.Errorf("%s rq.URL.Fragment = %v; want %v", identifier, got, want)
			}
		}
	}
}
