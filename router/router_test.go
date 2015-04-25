package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeHandler struct {
	called bool
}

func (f *fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.called = true
}

func TestMountPath(t *testing.T) {
	last := New()
	middle := New()
	middle2 := New()
	parent := New()
	fake := New()
	fake.prefix = "f"

	last.Mount("b", middle2)
	middle2.Mount("y", middle)
	middle.Mount("a", parent)
	fh := &fakeHandler{}
	last.INDEX("x", fh)

	if mp := parent.MountPath(); mp != "/" {
		t.Errorf("parent mountpath not /, but: %#v", mp)
	}

	if mp := middle.MountPath(); mp != "/a" {
		t.Errorf("middle mountpath not /a, but: %#v", mp)
	}

	if mp := last.MountPath(); mp != "/a/y/b" {
		t.Errorf("last mountpath not /a/y/b, but: %#v", mp)
	}

	if mp := fake.MountPath(); mp != "f" {
		t.Errorf("fake mountpath not f, but: %#v", mp)
	}

	rq, _ := http.NewRequest("GET", "/a/y/b/x", nil)
	rec := httptest.NewRecorder()

	parent.ServeHTTP(rec, rq)

	if !fh.called {
		t.Errorf("not called")
	}
}

func TestMountRouterTwice(t *testing.T) {
	rt := New()
	parent := New()
	rt.Mount("r", parent)

	mustPanic(t, func() {
		rt.Mount("r", parent)
	})
}

func TestInvalidSegment(t *testing.T) {
	rt := New()
	mustPanic(t, func() {
		rt.POST("a/b", &fakeHandler{})
	})

	mustPanic(t, func() {
		rt.INDEX("a/b", &fakeHandler{})
	})
}

func TestIndex(t *testing.T) {
	rt := New()
	fh := &fakeHandler{}
	url := rt.INDEXFunc("", fh.ServeHTTP)

	if url() != "/" {
		t.Errorf("wrong url: %s // expected: /", url())
	}

	rec := httptest.NewRecorder()
	rq, _ := http.NewRequest("GET", "/", nil)

	rt.ServeHTTP(rec, rq)

	if !fh.called {
		t.Errorf("handler not called")
	}
}

func TestPost(t *testing.T) {
	rt := New()
	fh := &fakeHandler{}
	url := rt.POSTFunc("a", fh.ServeHTTP)

	if url() != "/a" {
		t.Errorf("wrong url: %s // expected: /a", url())
	}

	rec := httptest.NewRecorder()
	rq, _ := http.NewRequest("POST", "/a", nil)

	rt.ServeHTTP(rec, rq)

	if !fh.called {
		t.Errorf("handler not called")
	}
}

func TestIndexMounted(t *testing.T) {
	rt := New()
	fh := &fakeHandler{}
	url := rt.INDEXFunc("", fh.ServeHTTP)
	parent := New()
	rt.Mount("c", parent)

	if url() != "/c" {
		t.Errorf("wrong url: %s // expected: /c", url())
	}

	rec := httptest.NewRecorder()
	rq, _ := http.NewRequest("GET", "/c", nil)

	parent.ServeHTTP(rec, rq)

	if !fh.called {
		t.Errorf("handler not called")
	}
}

func TestMultiRoutes(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.GET("a", h)
	rt.POST("a", h)
	rt.PUT("a", h)
	rt.DELETE("a", h)
	rt.PATCH("a", h)

	rt.GET("b", h)
	rt.POST("b", h)
	rt.PUT("b", h)
	rt.DELETE("b", h)
	rt.PATCH("b", h)
}

func TestDoubleRouteHandler1(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.Handle("a", h)

	mustPanic(t, func() {
		rt.INDEX("a", h)
	})

	mustPanic(t, func() {
		rt.POST("a", h)
	})
}

func TestDoubleHandler1(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.Handle("a", h)

	mustPanic(t, func() {
		rt.Handle("a", h)
	})

	mustPanic(t, func() {
		New().Mount("a", rt)
	})
}

func TestRegisterInvalidMethod(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	mustPanic(t, func() {
		rt.registerHandler("x", "y", h)
	})
}

func TestInvalidSegmentForHandler(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	mustPanic(t, func() {
		rt.Handle("a/b", h)
	})
}

func TestDoubleHandler2(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.GET("a", h)
	rt.POST("b", h)
	rt.PATCH("c", h)
	rt.PUT("d", h)
	rt.DELETE("e", h)
	rt.INDEX("f", h)

	mustPanic(t, func() {
		rt.Handle("a", h)
	})

	mustPanic(t, func() {
		rt.Handle("b", h)
	})

	mustPanic(t, func() {
		rt.Handle("c", h)
	})

	mustPanic(t, func() {
		rt.Handle("d", h)
	})

	mustPanic(t, func() {
		rt.Handle("e", h)
	})
	mustPanic(t, func() {
		rt.Handle("f", h)
	})
}

func TestDoubleRouteHandler2(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.Handle("a", h)

	mustPanic(t, func() {
		rt.GET("a", h)
	})
}

func TestDoubleRouteGET(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.GET("a", h)

	mustPanic(t, func() {
		rt.GET("a", h)
	})
}

func TestDoubleRouteDELETE(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.DELETE("a", h)

	mustPanic(t, func() {
		rt.DELETE("a", h)
	})
}

func TestDoubleRoutePUT(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.PUT("a", h)

	mustPanic(t, func() {
		rt.PUT("a", h)
	})
}

func TestDoubleRoutePOST(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.POST("a", h)

	mustPanic(t, func() {
		rt.POST("a", h)
	})
}

func TestDoubleRoutePATCH(t *testing.T) {
	rt := New()
	h := &fakeHandler{}
	rt.PATCH("a", h)

	mustPanic(t, func() {
		rt.PATCH("a", h)
	})
}

func TestGet(t *testing.T) {
	rt := New()
	fh := &fakeHandler{}
	url := rt.GET("a", fh)

	if url("3") != "/a/3" {
		t.Errorf("wrong url: %s // expected: /a/3", url("3"))
	}
}

func TestGetMounted(t *testing.T) {
	rt := New()
	fh := &fakeHandler{}
	url := rt.GET("a", fh)
	parent := New()
	rt.Mount("c", parent)

	if url("3") != "/c/a/3" {
		t.Errorf("wrong url: %s // expected: /c/a/3", url("3"))
	}
}

func TestSearch(t *testing.T) {
	rt := New()
	fh := &fakeHandler{}
	url := rt.INDEX("a", fh)
	parent := New()
	rt.Mount("c", parent)

	if url() != "/c/a" {
		t.Errorf("wrong url: %s // expected: /c/a", url())
	}
}

func TestInvalidSegmentParams(t *testing.T) {
	rt := New()
	mustPanic(t, func() {
		rt.GET("a/b", &fakeHandler{})
	})
}

func TestHandleRouterPanic(t *testing.T) {
	parent := New()
	rt := New()

	mustPanic(t, func() {
		parent.Handle("r", rt)
	})
}

func TestHandleFunc(t *testing.T) {
	// t.SkipNow()
	tests := []struct {
		requestMethod string
		requestPath   string
		pathSegment   string
		called        bool
		code          int
		mountPoint    string
	}{
		{"GET", "/a", "a", true, 200, ""},
		{"GET", "/a/", "a", true, 200, ""},
		{"GET", "/a/1", "a", true, 200, ""},
		{"GET", "/a/b/1", "a", true, 200, ""},
		{"GET", "/b/1", "a", false, 404, ""},

		{"POST", "/b", "b", true, 200, ""},
		{"POST", "/b/", "b", true, 200, ""},
		{"POST", "/b/1", "b", true, 200, ""},
		{"POST", "/b/b/1", "b", true, 200, ""},
		{"POST", "/b/1", "a", false, 405, ""},

		{"PUT", "/c", "c", true, 200, ""},
		{"PUT", "/c/", "c", true, 200, ""},
		{"PUT", "/c/1", "c", true, 200, ""},
		{"PUT", "/c/b/1", "c", true, 200, ""},
		{"PUT", "/b/1", "a", false, 405, ""},

		{"PATCH", "/d", "d", true, 200, ""},
		{"PATCH", "/d/", "d", true, 200, ""},
		{"PATCH", "/d/1", "d", true, 200, ""},
		{"PATCH", "/d/b/1", "d", true, 200, ""},
		{"PATCH", "/b/1", "a", false, 405, ""},

		{"DELETE", "/e", "e", true, 200, ""},
		{"DELETE", "/e/", "e", true, 200, ""},
		{"DELETE", "/e/1", "e", true, 200, ""},
		{"DELETE", "/e/b/1", "e", true, 200, ""},
		{"DELETE", "/b/1", "a", false, 405, ""},

		// Mounted

		{"GET", "/m/a", "a", true, 200, "m"},
		{"GET", "/m/a/", "a", true, 200, "m"},
		{"GET", "/m/a/1", "a", true, 200, "m"},
		{"GET", "/m/a/b/1", "a", true, 200, "m"},
		{"GET", "/m/b/1", "a", false, 404, "m"},
		{"GET", "/x/a", "a", false, 404, "m"},
		{"GET", "/a", "a", false, 404, "m"},

		{"POST", "/m/b", "b", true, 200, "m"},
		{"POST", "/m/b/", "b", true, 200, "m"},
		{"POST", "/m/b/1", "b", true, 200, "m"},
		{"POST", "/m/b/b/1", "b", true, 200, "m"},
		{"POST", "/m/a/1", "b", false, 405, "m"},
		{"POST", "/x/b", "b", false, 405, "m"},
		{"POST", "/b", "b", false, 405, "m"},

		{"PUT", "/m/c", "c", true, 200, "m"},
		{"PUT", "/m/c/", "c", true, 200, "m"},
		{"PUT", "/m/c/1", "c", true, 200, "m"},
		{"PUT", "/m/c/b/1", "c", true, 200, "m"},
		{"PUT", "/m/b/1", "c", false, 405, "m"},
		{"PUT", "/x/c", "c", false, 405, "m"},
		{"PUT", "/c", "c", false, 405, "m"},

		{"PATCH", "/m/d", "d", true, 200, "m"},
		{"PATCH", "/m/d/", "d", true, 200, "m"},
		{"PATCH", "/m/d/1", "d", true, 200, "m"},
		{"PATCH", "/m/d/b/1", "d", true, 200, "m"},
		{"PATCH", "/m/b/1", "d", false, 405, "m"},
		{"PATCH", "/x/d", "d", false, 405, "m"},
		{"PATCH", "/d", "d", false, 405, "m"},

		{"DELETE", "/m/e", "e", true, 200, "m"},
		{"DELETE", "/m/e/", "e", true, 200, "m"},
		{"DELETE", "/m/e/1", "e", true, 200, "m"},
		{"DELETE", "/m/e/b/1", "e", true, 200, "m"},
		{"DELETE", "/m/b/1", "e", false, 405, "m"},
		{"DELETE", "/x/e", "e", false, 405, "m"},
		{"DELETE", "/e", "e", false, 405, "m"},
	}

	for _, test := range tests {

		rt := New()
		fh := &fakeHandler{}
		urlStr := test.requestPath

		rt.HandleFunc(test.pathSegment, fh.ServeHTTP)
		rq, err := http.NewRequest(test.requestMethod, urlStr, nil)
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

		identifier := test.requestMethod + " " + urlStr + " for " + test.pathSegment

		if test.mountPoint != "" {
			identifier += " mounted at " + test.mountPoint
		}

		if got, want := rec.Code, test.code; got != want {
			t.Errorf("%s rec.Code = %v; want %v", identifier, got, want)
		}

		if got, want := fh.called, test.called; got != want {
			t.Errorf("%s fh.called = %v; want %v", identifier, got, want)
		}
	}

}

func TestMethodHandlers(t *testing.T) {
	// t.SkipNow()
	tests := []struct {
		requestMethod string
		requestPath   string
		method        string
		pathSegment   string
		param         string
		called        bool
		code          int
		mountPoint    string
	}{
		// GET for details
		{"GET", "/a/1", "GET", "a", "1", true, 200, ""},     // working
		{"GET", "/a/1/", "GET", "a", "1", true, 200, ""},    // working: additional slash
		{"GET", "/a", "GET", "a", "1", false, 404, ""},      // wrong GET: parameter missing (1)
		{"GET", "/a/", "GET", "a", "1", false, 404, ""},     // wrong GET: parameter missing (2)
		{"GET", "/a/1", "GET", "b", "1", false, 404, ""},    // wrong url
		{"POST", "/a/1", "GET", "a", "1", false, 405, ""},   // wrong method: POST instead of GET
		{"PUT", "/a/1", "GET", "a", "1", false, 405, ""},    // wrong method: PUT instead of GET
		{"PATCH", "/a/1", "GET", "a", "1", false, 405, ""},  // wrong method: PATCH instead of GET
		{"DELETE", "/a/1", "GET", "a", "1", false, 405, ""}, // wrong method: DELETE instead of GET

		// GET without details
		{"GET", "/b", "GET", "b", "", true, 200, ""},     // working
		{"GET", "/b/", "GET", "b", "", true, 200, ""},    // working: additional slash
		{"GET", "/b/1", "GET", "b", "", false, 404, ""},  // wrong GET: additional parameter
		{"GET", "/b", "GET", "c", "", false, 404, ""},    // wrong url
		{"POST", "/b", "GET", "b", "", false, 405, ""},   // wrong method: POST instead of GET
		{"PUT", "/b", "GET", "b", "", false, 405, ""},    // wrong method: PUT instead of GET
		{"PATCH", "/b", "GET", "b", "", false, 405, ""},  // wrong method: PATCH instead of GET
		{"DELETE", "/b", "GET", "b", "", false, 405, ""}, // wrong method: DELETE instead of GET

		// POST
		{"POST", "/c", "POST", "c", "", true, 200, ""},    // working
		{"POST", "/c/", "POST", "c", "", true, 200, ""},   // working: additional slash
		{"POST", "/c/1", "POST", "c", "", false, 405, ""}, // wrong GET: additional parameter
		{"POST", "/c", "POST", "d", "", false, 405, ""},   // wrong url
		{"GET", "/c", "POST", "c", "", false, 404, ""},    // wrong method: GET instead of POST
		{"PUT", "/c", "POST", "c", "", false, 405, ""},    // wrong method: PUT instead of POST
		{"PATCH", "/c", "POST", "c", "", false, 405, ""},  // wrong method: PATCH instead of POST
		{"DELETE", "/c", "POST", "c", "", false, 405, ""}, // wrong method: DELETE instead of POST

		// PUT
		{"PUT", "/d/1", "PUT", "d", "1", true, 200, ""},     // working
		{"PUT", "/d/1/", "PUT", "d", "1", true, 200, ""},    // working: additional slash
		{"PUT", "/d", "PUT", "d", "1", false, 405, ""},      // wrong PUT: parameter missing (1)
		{"PUT", "/d/", "PUT", "d", "1", false, 405, ""},     // wrong PUT: parameter missing (2)
		{"PUT", "/d/1", "PUT", "b", "1", false, 405, ""},    // wrong URL
		{"POST", "/d/1", "PUT", "d", "1", false, 405, ""},   // wrong method: POST instead of PUT
		{"GET", "/d/1", "PUT", "d", "1", false, 404, ""},    // wrong method: GET instead of PUT
		{"PATCH", "/d/1", "PUT", "d", "1", false, 405, ""},  // wrong method: PATCH instead of PUT
		{"DELETE", "/d/1", "PUT", "d", "1", false, 405, ""}, // wrong method: DELETE instead of PUT

		// PATCH
		{"PATCH", "/e/1", "PATCH", "e", "1", true, 200, ""},   // working
		{"PATCH", "/e/1/", "PATCH", "e", "1", true, 200, ""},  // working: additional slash
		{"PATCH", "/e", "PATCH", "e", "1", false, 405, ""},    // wrong PATCH: parameter missing (1)
		{"PATCH", "/e/", "PATCH", "e", "1", false, 405, ""},   // wrong PATCH: parameter missing (2)
		{"PATCH", "/e/1", "PATCH", "b", "1", false, 405, ""},  // wrong URL
		{"POST", "/e/1", "PATCH", "e", "1", false, 405, ""},   // wrong method: POST instead of PATCH
		{"PUT", "/e/1", "PATCH", "e", "1", false, 405, ""},    // wrong method: PUT instead of PATCH
		{"GET", "/e/1", "PATCH", "e", "1", false, 404, ""},    // wrong method: GET instead of PATCH
		{"DELETE", "/e/1", "PATCH", "e", "1", false, 405, ""}, // wrong method: DELETE instead of PATCH

		// DELETE
		{"DELETE", "/f/1", "DELETE", "f", "1", true, 200, ""},  // working
		{"DELETE", "/f/1/", "DELETE", "f", "1", true, 200, ""}, // working: additional slash
		{"DELETE", "/f", "DELETE", "f", "1", false, 405, ""},   // wrong DELETE: parameter missing (1)
		{"DELETE", "/f/", "DELETE", "f", "1", false, 405, ""},  // wrong DELETE: parameter missing (2)
		{"DELETE", "/f/1", "DELETE", "b", "1", false, 405, ""}, // wrong URL
		{"POST", "/f/1", "DELETE", "f", "1", false, 405, ""},   // wrong method: POST instead of DELETE
		{"PUT", "/f/1", "DELETE", "f", "1", false, 405, ""},    // wrong method: PUT instead of DELETE
		{"PATCH", "/f/1", "DELETE", "f", "1", false, 405, ""},  // wrong method: PATCH instead of DELETE
		{"GET", "/f/1", "DELETE", "f", "1", false, 404, ""},    // wrong method: GET instead of DELETE

		// Mounted

		// GET mounted for details
		{"GET", "/m/a/1", "GET", "a", "1", true, 200, "m"},     // working
		{"GET", "/m/a/1/", "GET", "a", "1", true, 200, "m"},    // working: additional slash
		{"GET", "/m/a", "GET", "a", "1", false, 404, "m"},      // wrong GET: parameter missing (1)
		{"GET", "/m/a/", "GET", "a", "1", false, 404, "m"},     // wrong GET: parameter missing (2)
		{"GET", "/m/a/1", "GET", "b", "1", false, 404, "m"},    // wrong url
		{"POST", "/m/a/1", "GET", "a", "1", false, 405, "m"},   // wrong method: POST instead of GET
		{"PUT", "/m/a/1", "GET", "a", "1", false, 405, "m"},    // wrong method: PUT instead of GET
		{"PATCH", "/m/a/1", "GET", "a", "1", false, 405, "m"},  // wrong method: PATCH instead of GET
		{"DELETE", "/m/a/1", "GET", "a", "1", false, 405, "m"}, // wrong method: DELETE instead of GET
		{"GET", "/x/a/1", "GET", "a", "1", false, 404, "m"},    // wrong mountpoint
		{"GET", "/a/1", "GET", "a", "1", false, 404, "m"},      // missing mountpoint

		// GET mounted without details
		{"GET", "/m/b", "GET", "b", "", true, 200, "m"},     // working
		{"GET", "/m/b/", "GET", "b", "", true, 200, "m"},    // working: additional slash
		{"GET", "/m/b/1", "GET", "b", "", false, 404, "m"},  // wrong GET: additional parameter
		{"GET", "/m/b", "GET", "c", "", false, 404, "m"},    // wrong url
		{"POST", "/m/b", "GET", "b", "", false, 405, "m"},   // wrong method: POST instead of GET
		{"PUT", "/m/b", "GET", "b", "", false, 405, "m"},    // wrong method: PUT instead of GET
		{"PATCH", "/m/b", "GET", "b", "", false, 405, "m"},  // wrong method: PATCH instead of GET
		{"DELETE", "/m/b", "GET", "b", "", false, 405, "m"}, // wrong method: DELETE instead of GET
		{"GET", "/x/b", "GET", "b", "", false, 404, "m"},    // wrong mountpoint
		{"GET", "/b", "GET", "b", "1", false, 404, "m"},     // missing mountpoint

		// POST mounted
		{"POST", "/m/c", "POST", "c", "", true, 200, "m"},    // working
		{"POST", "/m/c/", "POST", "c", "", true, 200, "m"},   // working: additional slash
		{"POST", "/m/c/1", "POST", "c", "", false, 405, "m"}, // wrong GET: additional parameter
		{"POST", "/m/c", "POST", "d", "", false, 405, "m"},   // wrong url
		{"GET", "/m/c", "POST", "c", "", false, 404, "m"},    // wrong method: GET instead of POST
		{"PUT", "/m/c", "POST", "c", "", false, 405, "m"},    // wrong method: PUT instead of POST
		{"PATCH", "/m/c", "POST", "c", "", false, 405, "m"},  // wrong method: PATCH instead of POST
		{"DELETE", "/m/c", "POST", "c", "", false, 405, "m"}, // wrong method: DELETE instead of POST
		{"POST", "/x/c", "POST", "c", "", false, 405, "m"},   // wrong mountpoint
		{"POST", "/c", "POST", "c", "", false, 405, "m"},     // missing mountpoint

		// PUT mounted
		{"PUT", "/m/d/1", "PUT", "d", "1", true, 200, "m"},     // working
		{"PUT", "/m/d/1/", "PUT", "d", "1", true, 200, "m"},    // working: additional slash
		{"PUT", "/m/d", "PUT", "d", "1", false, 405, "m"},      // wrong PUT: parameter missing (1)
		{"PUT", "/m/d/", "PUT", "d", "1", false, 405, "m"},     // wrong PUT: parameter missing (2)
		{"PUT", "/m/d/1", "PUT", "b", "1", false, 405, "m"},    // wrong URL
		{"POST", "/m/d/1", "PUT", "d", "1", false, 405, "m"},   // wrong method: POST instead of PUT
		{"GET", "/m/d/1", "PUT", "d", "1", false, 404, "m"},    // wrong method: GET instead of PUT
		{"PATCH", "/m/d/1", "PUT", "d", "1", false, 405, "m"},  // wrong method: PATCH instead of PUT
		{"DELETE", "/m/d/1", "PUT", "d", "1", false, 405, "m"}, // wrong method: DELETE instead of PUT
		{"PUT", "/x/d/1", "PUT", "d", "1", false, 405, "m"},    // wrong mountpoint
		{"PUT", "/d/1", "PUT", "d", "1", false, 405, "m"},      // missing mountpoint

		// PATCH mounted
		{"PATCH", "/m/e/1", "PATCH", "e", "1", true, 200, "m"},   // working
		{"PATCH", "/m/e/1/", "PATCH", "e", "1", true, 200, "m"},  // working: additional slash
		{"PATCH", "/m/e", "PATCH", "e", "1", false, 405, "m"},    // wrong PATCH: parameter missing (1)
		{"PATCH", "/m/e/", "PATCH", "e", "1", false, 405, "m"},   // wrong PATCH: parameter missing (2)
		{"PATCH", "/m/e/1", "PATCH", "b", "1", false, 405, "m"},  // wrong URL
		{"POST", "/m/e/1", "PATCH", "e", "1", false, 405, "m"},   // wrong method: POST instead of PATCH
		{"PUT", "/m/e/1", "PATCH", "e", "1", false, 405, "m"},    // wrong method: PUT instead of PATCH
		{"GET", "/m/e/1", "PATCH", "e", "1", false, 404, "m"},    // wrong method: GET instead of PATCH
		{"DELETE", "/m/e/1", "PATCH", "e", "1", false, 405, "m"}, // wrong method: DELETE instead of PATCH
		{"PATCH", "/x/e/1", "PATCH", "e", "1", false, 405, "m"},  // wrong mountpoint
		{"PATCH", "/e/1", "PATCH", "e", "1", false, 405, "m"},    // missing mountpoint

		// DELETE mounted
		{"DELETE", "/m/f/1", "DELETE", "f", "1", true, 200, "m"},  // working
		{"DELETE", "/m/f/1/", "DELETE", "f", "1", true, 200, "m"}, // working: additional slash
		{"DELETE", "/m/f", "DELETE", "f", "1", false, 405, "m"},   // wrong DELETE: parameter missing (1)
		{"DELETE", "/m/f/", "DELETE", "f", "1", false, 405, "m"},  // wrong DELETE: parameter missing (2)
		{"DELETE", "/m/f/1", "DELETE", "b", "1", false, 405, "m"}, // wrong URL
		{"POST", "/m/f/1", "DELETE", "f", "1", false, 405, "m"},   // wrong method: POST instead of DELETE
		{"PUT", "/m/f/1", "DELETE", "f", "1", false, 405, "m"},    // wrong method: PUT instead of DELETE
		{"PATCH", "/m/f/1", "DELETE", "f", "1", false, 405, "m"},  // wrong method: PATCH instead of DELETE
		{"GET", "/m/f/1", "DELETE", "f", "1", false, 404, "m"},    // wrong method: GET instead of DELETE
		{"DELETE", "/x/f/1", "DELETE", "f", "1", false, 405, "m"}, // wrong mountpoint
		{"DELETE", "/f/1", "DELETE", "f", "1", false, 405, "m"},   // missing mountpoint

	}

	for _, test := range tests {
		rt := New()
		fh := &fakeHandler{}
		urlStr := test.requestPath

		if test.param != "" {
			rt.routeParam(test.method, test.pathSegment, fh)
		} else {
			rt.route(test.method, test.pathSegment, fh)
		}

		rq, err := http.NewRequest(test.requestMethod, urlStr, nil)
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

		identifier := test.requestMethod + " " + urlStr + " for " + test.method + " " + test.pathSegment

		if test.mountPoint != "" {
			identifier += " mounted at " + test.mountPoint
		}

		if got, want := rec.Code, test.code; got != want {
			t.Errorf("%s rec.Code = %v; want %v", identifier, got, want)
		}

		if got, want := fh.called, test.called; got != want {
			t.Errorf("%s fh.called = %v; want %v", identifier, got, want)
		}

		if test.param != "" && fh.called {
			if got, want := rq.URL.Fragment, test.param; got != want {
				t.Errorf("%s rq.URL.Fragment = %v; want %v", identifier, got, want)
			}
		}
	}
}
