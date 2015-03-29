package main

import (
	"fmt"
	"github.com/go-on/stack/mw"
	"github.com/go-on/stack/router"
	"net/http"
)

var (
	urls struct {
		a, idx, z func() string
		hello     func(string) string
	}

	layoutMW = []interface{}{
		mw.Around(mw.HTMLString("<h1>Example</h1>"), mw.HTMLString("<h2>End</h2>")),
	}

	middle = router.New()
	sub    = router.New()
)

func init() {
	urls.a = middle.INDEX("a", mw.String("A"))
	urls.hello = middle.GETFunc("hello", hello)
	urls.idx = middle.INDEXFunc("", index)
	urls.z = sub.INDEX("z", mw.String("Z"))
}

func main() {
	top := router.New()
	sub.Mount("p", middle)
	middle.Mount("x", top, layoutMW...)

	top.INDEX("", mw.TemporaryRedirect(urls.idx()))
	http.ListenAndServe(":8080", top)
}

func hello(rw http.ResponseWriter, req *http.Request) {
	mw.HTMLString(fmt.Sprintf(`
		<p>
			Hello, %s!
		</p>
	`,
		req.URL.Fragment,
	)).ServeHTTP(rw, req)
}

func index(rw http.ResponseWriter, req *http.Request) {
	mw.HTMLString(fmt.Sprintf(`
		<ul>
			<li><a href="%s">a</a></li>
			<li><a href="%s">hello peter</a></li>
			<li><a href="%s">z</a></li>
			<li><a href="/p">nothing</a></li>
		</ul>
	`,
		urls.a(),
		urls.hello("Peter"),
		urls.z(),
	)).ServeHTTP(rw, req)
}
