package mux

import (
	"fmt"
	"gopkg.in/go-on/stack.v4"
	"net/http"
	"strings"
	"sync"
)

func New() *Mux {
	return &Mux{edges: make(map[string]*Mux)}
}

type Mux struct {
	edges   map[string]*Mux // the empty key is for the next wildcard Mux (the Mux after my wildcard)
	handler http.Handler
	Debug   bool
}

func (p *Mux) HandleFunc(path string, fn func(http.ResponseWriter, *http.Request)) {
	p.Handle(path, http.HandlerFunc(fn))
}

func (p *Mux) Handle(path string, rh http.Handler) {
	if p.Debug {
		fmt.Printf("registering handler %T for path %#v\n", rh, path)
	}
	/*
		if path[0] != '/' || path[len(path)-1] != '/' {
			panic("path must start and end with /")
		}
	*/
	nd := p
	var start int = 1
	var end int
	var fin bool
	for {
		if start >= len(path) {
			break
		}
		end = strings.Index(path[start:], "/")

		if end == 0 {
			start++
			continue
		}

		if end == -1 {
			end = len(path) - 1
			fin = true
		} else {
			end += start
		}

		pa := path[start : end+1]
		// println("registering", pa)
		subMux := nd.edges[pa]
		if subMux == nil {
			subMux = &Mux{edges: make(map[string]*Mux)}
			nd.edges[pa] = subMux
		}
		nd = subMux

		if fin {
			break
		}

		start = end + 1
	}
	nd.handler = rh
}

func (m *Mux) findEdge(path string) *Mux {
	// println("findEdge", path)
	if len(m.edges) > 0 {
		hd, has := m.edges[path]
		if has {
			// println("found")
			return hd
		}
	}
	// println("not found")
	return nil
}

func (m *Mux) findHandler(start int, req *http.Request) http.Handler {
	endPath := len(req.URL.Path)
	if start >= endPath {
		return m.handler
	}

	edge := m.findEdge(req.URL.Path[start:])

	if edge != nil {
		return edge.handler
	}

	pos := strings.Index(req.URL.Path[start:], "/")
	end := start + pos
	if pos == -1 {
		return m.handler
	}

	pth := req.URL.Path[start : end+1]

	edge = m.findEdge(pth)

	if edge != nil {
		return edge.findHandler(end+1, req)
	}

	return m.handler
}

func (n *Mux) Handler(req *http.Request) http.Handler {
	return n.findHandler(1, req)
}

func (n *Mux) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	h := n.Handler(req)

	if h == nil {
		http.NotFound(wr, req)
		return
	}

	h.ServeHTTP(wr, req)
}

func (m *Mux) Mount(mountPoint string, ma Mountable) error {
	return ma.Mount(m, mountPoint)
}

func (m *Mux) MountWrapped(mountPoint string, ma Mountable, wr stack.Wrapper) error {
	wrmx := NewWrappedMux(wr)
	err := ma.Mount(wrmx, mountPoint)
	if err != nil {
		panic(err)
	}
	m.Handle(mountPoint, wrmx)
	return nil
}

func (m *Mux) MustMount(mountPoint string, ma Mountable) {
	err := m.Mount(mountPoint, ma)
	if err != nil {
		panic(err)
	}
}

func (m *Mux) MustMountWrapped(mountPoint string, ma Mountable, wr stack.Wrapper) {
	err := m.MountWrapped(mountPoint, ma, wr)
	if err != nil {
		panic(err)
	}
}

// Mount mounts on the http.DefaultServeMux
func Mount(mountPoint string, ma Mountable) error {
	return ma.Mount(http.DefaultServeMux, mountPoint)
}

// MountWrapps mounts on the http.DefaultServeMux
func MountWrapped(mountPoint string, ma Mountable, wr stack.Wrapper) error {
	wrmx := NewWrappedMux(wr)
	err := ma.Mount(wrmx, mountPoint)
	if err != nil {
		panic(err)
	}
	http.Handle(mountPoint, wrmx)
	return nil
}

func MustMount(mountPoint string, ma Mountable) {
	err := Mount(mountPoint, ma)
	if err != nil {
		panic(err)
	}
}

func MustMountWrapped(mountPoint string, ma Mountable, wr stack.Wrapper) {
	err := MountWrapped(mountPoint, ma, wr)
	if err != nil {
		panic(err)
	}
}

type Mountable interface {
	Mount(mx Muxer, mountPoint string) error
}

// WrappedMux is a Muxer that wraps each handler that is registered by a
// wrapper
type WrappedMux struct {
	mx *Mux
	wr stack.Wrapper
	sync.Once
	h http.Handler
}

// Handle calls the Handle method of the muxer and passes the wrapped handler
func (wm *WrappedMux) Handle(path string, h http.Handler) {
	wm.mx.Handle(path, h)
}

func (wm *WrappedMux) HandleFunc(path string, fn func(http.ResponseWriter, *http.Request)) {
	wm.Handle(path, http.HandlerFunc(fn))
}

func (wm *WrappedMux) Debug() {
	wm.mx.Debug = true
}

func (wm *WrappedMux) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	wm.Once.Do(func() {
		wm.h = wm.wr.Wrap(wm.mx)
	})
	wm.h.ServeHTTP(wr, req)
}

func NewWrappedMux(wr stack.Wrapper) *WrappedMux {
	return &WrappedMux{mx: New(), wr: wr}
}

type Muxer interface {
	Handle(path string, h http.Handler)
	HandleFunc(path string, fn func(http.ResponseWriter, *http.Request))
	ServeHTTP(wr http.ResponseWriter, req *http.Request)
}
