package stacksession_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	"github.com/go-on/stack"
	"github.com/go-on/stack/third-party/stacksession"
	"gopkg.in/go-on/sessions.v1"
)

func printFlashes(ctx stack.Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	var session stacksession.Session
	ctx.Get(&session)
	for _, fl := range session.Flashes() {
		fmt.Printf("Flash: %v\n", fl)
	}
	next.ServeHTTP(w, r)
}

var counter = int64(0)

func setNameAndFlash(ctx stack.Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
	atomic.AddInt64(&counter, 1)
	var session stacksession.Session
	ctx.Get(&session)
	session.Values["name"] = r.URL.Query().Get("name")
	if counter == 1 {
		session.AddFlash("Hello, flash messages world!")
	}
	next.ServeHTTP(w, r)
}

func printName(ctx stack.Contexter, w http.ResponseWriter, r *http.Request) {
	var session stacksession.Session
	ctx.Get(&session)
	fmt.Printf("Name: %s\n", session.Values["name"])
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func Example() {
	s := stack.New().
		UseWithContext(stacksession.SaveAndClear).
		UseWithContext(stacksession.NewStore(store, "my-session-name")).
		UseFuncWithContext(printFlashes).
		UseFuncWithContext(setNameAndFlash).
		WrapFuncWithContext(printName)

	req, _ := http.NewRequest("GET", "/?name=Peter", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	s.ServeHTTP(rec, req)
	s.ServeHTTP(rec, req)

	// Output:
	// Name: Peter
	// Flash: Hello, flash messages world!
	// Name: Peter
	// Name: Peter
}
