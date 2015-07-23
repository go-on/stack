package main

import (
	"fmt"
	"gopkg.in/go-on/stack.v1"
	"gopkg.in/go-on/stack.v1/mux"
	"gopkg.in/go-on/stack.v1/mw"
	"math/rand"
	"net/http"
	"time"
	// "gopkg.in/go-on/stack.v1/rest"
	"gopkg.in/go-on/stack.v1/rest/example/app"
	"gopkg.in/go-on/stack.v1/server"
)

var layoutMW = stack.New().Use(
	mw.Around(
		mw.HTMLString("<h1>Example</h1>"),
		mw.HTMLString("<h2>End</h2>"),
	),
)

func main1(a *app.App, mountPt string, layout bool) {
	// uses http.DefaultMux
	if layout {
		mux.MustMountWrapped(mountPt, a, layoutMW)
	} else {
		mux.MustMount(mountPt, a)
	}
	fmt.Printf("urls: %v\n", a.URLs)

	// uses http.DefaultMux
	http.ListenAndServe(":8080", nil)
}

func main2(a *app.App, mountPt string, layout bool) {
	mx := mux.New()
	if layout {
		mx.MustMountWrapped(mountPt, a, layoutMW)
	} else {
		mx.MustMount(mountPt, a)
	}
	fmt.Printf("urls: %v\n", a.URLs)

	http.ListenAndServe(":8080", mx)
}

func main3(a *app.App, mountPt string, layout bool) {
	s := server.New()
	if layout {
		s.MustMountWrapped(mountPt, a, layoutMW)
	} else {
		s.MustMount(mountPt, a)
	}
	fmt.Printf("urls: %v\n", a.URLs)

	s.HTTPServer(":8080").ListenAndServe()
}

func main() {

	// main1-3 do the same, pick one
	i := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(3)
	layout := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(2) == 1

	a := app.New()
	mountPt := fmt.Sprintf("/main%d/", i+1)
	fmt.Printf("running %s with layout: %v\n", mountPt, layout)
	switch i {
	case 0:
		main1(a, mountPt, layout)
	case 1:
		main2(a, mountPt, layout)
	case 2:
		main3(a, mountPt, layout)
	}

}
