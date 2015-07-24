package app

import (
	"fmt"
	"gopkg.in/go-on/stack.v6/mux"
	"gopkg.in/go-on/stack.v6/mw"
	"gopkg.in/go-on/stack.v6/rest"
	"net/http"
)

type App struct {
	URLs struct {
		A, Idx, Z, Hello string
	}
}

func New() *App {
	return &App{}
}

func (a *App) Mount(m mux.Muxer, mountpoint string) error {
	if m == nil {
		m = http.DefaultServeMux
	}

	a.setURLs(mountpoint)
	a.mountHandlers(m)
	return nil
}

func (a *App) setURLs(mountpoint string) {
	a.URLs.Idx = mountpoint
	a.URLs.A = mountpoint + "a/"
	a.URLs.Hello = mountpoint + "hello/"
	a.URLs.Z = mountpoint + "p/z/"
}

func (a *App) mountHandlers(m mux.Muxer) {
	aMux := rest.
		Index(mw.String("A")).
		GetFunc(a.getA)

	m.Handle(a.URLs.A, aMux)
	m.Handle(a.URLs.Idx, rest.IndexFunc(a.index))
	m.Handle(a.URLs.Hello, rest.GetFunc(a.hello))
	m.Handle(a.URLs.Z, rest.Index(mw.String("Z")))

}

func (a *App) getA(rw http.ResponseWriter, req *http.Request) {
	mw.HTMLString(fmt.Sprintf(`
		<p>
			A with %s!
		</p>
	`,
		rest.Param(req),
	)).ServeHTTP(rw, req)
}

func (a *App) hello(rw http.ResponseWriter, req *http.Request) {
	mw.HTMLString(fmt.Sprintf(`
		<p>
			Hello, %s!
		</p>
	`,
		rest.Param(req),
	)).ServeHTTP(rw, req)
}

func (a *App) index(rw http.ResponseWriter, req *http.Request) {
	mw.HTMLString(fmt.Sprintf(`
		<ul>
			<li><a href="%s">a</a></li>
			<li><a href="%s">a with B</a></li>
			<li><a href="%s">hello peter</a></li>
			<li><a href="%s">z</a></li>
			<li><a href="/p">nothing</a></li>
		</ul>
	`,
		a.URLs.A,
		a.URLs.A+"B",
		a.URLs.Hello+"Peter",
		a.URLs.Z,
	)).ServeHTTP(rw, req)
}
