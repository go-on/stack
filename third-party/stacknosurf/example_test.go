package stacknosurf_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/go-on/stack"
	"github.com/go-on/stack/third-party/stacknosurf"
)

// app serves the form value "a" for POST requests and otherwise the token
func app(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	if req.Method == "POST" {
		req.ParseForm()
		rw.Write([]byte(req.FormValue("a")))
		return
	}
	var token stacknosurf.Token

	ctx.Get(&token)
	rw.Write([]byte(string(token)))
}

func Example() {

	s := stack.New(
		stack.Context,
		&stacknosurf.CheckToken{},
		stacknosurf.SetToken{},
		app,
	)

	// here comes the tests
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	s.ServeHTTP(rec, req)
	token := rec.Body.String()
	cookie := parseCookie(rec)

	rec = httptest.NewRecorder()
	req = mkPostReq(cookie, token)
	s.ServeHTTP(rec, req)
	fmt.Println("-- success --")
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/", nil)
	s.ServeHTTP(rec, req)
	fmt.Println("-- fail --")
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())
	// Output:
	// -- success --
	// 200
	// b
	// -- fail --
	// 400
	//
}

func parseCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	cookie := rec.Header().Get("Set-Cookie")
	cookie2 := cookie[0:strings.Index(cookie, ";")]
	splitter := strings.Index(cookie2, "=")
	c := http.Cookie{}
	c.Name = cookie2[0:splitter]
	c.Value = cookie2[splitter+1:]
	return &c
}

func mkPostReq(cookie *http.Cookie, token string) *http.Request {
	var vals url.Values = map[string][]string{}
	vals.Set("a", "b")
	req, _ := http.NewRequest("POST", "http://localhost/", strings.NewReader(vals.Encode()))
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Referer", "http://localhost/")
	return req
}
