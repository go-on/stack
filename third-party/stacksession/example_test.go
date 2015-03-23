package stacksession_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-on/stack"
	"github.com/go-on/stack/third-party/stacksession"
	"gopkg.in/go-on/sessions.v1"
)

func sessionWrite(ctx stack.Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	var session stacksession.Session
	ctx.Get(&session)
	session.Values["name"] = r.URL.Query().Get("name")
	session.AddFlash("Hello, flash messages world!")
	next.ServeHTTP(w, r)
}

func sessionRead(ctx stack.Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	var session stacksession.Session
	ctx.Get(&session)
	fmt.Printf("Name: %s\n", session.Values["name"])
	for _, fl := range session.Flashes() {
		fmt.Printf("Flash: %v\n", fl)
	}
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func Example() {
	s := stack.New(
		stacksession.SaveAndClear,
		stacksession.NewStore(store, "my-session-name"),
		sessionWrite,
		sessionRead,
	).ContextHandler()

	req, _ := http.NewRequest("GET", "/?name=Peter", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	// Output:
	// Name: Peter
	// Flash: Hello, flash messages world!
}
