package mux

import (
	"github.com/go-on/stack"
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
}

func (p *Mux) HandleFunc(path string, fn func(http.ResponseWriter, *http.Request)) {
	p.Handle(path, http.HandlerFunc(fn))
}

func (p *Mux) Handle(path string, rh http.Handler) {
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
			end = len(path)
			fin = true
		} else {
			end += start
		}

		p := path[start:end]
		subMux := nd.edges[p]
		if subMux == nil {
			subMux = &Mux{edges: make(map[string]*Mux)}
			nd.edges[p] = subMux
		}
		nd = subMux

		if fin {
			break
		}

		start = end + 1
	}
	nd.handler = rh
}

func (n *Mux) findHandler(start int, endPath int, req *http.Request) http.Handler {
	if start > endPath {
		return n.handler
	}
	pos := strings.Index(req.URL.Path[start:endPath], "/")
	end := start + pos
	if pos == -1 {
		end = endPath
	}

	pth := req.URL.Path[start:end]
	for k, val := range n.edges {
		if k == pth {
			if len(val.edges) == 0 {
				return val.handler
			}
			return val.findHandler(end+1, endPath, req)
		}
	}
	return nil
}

func (n *Mux) Handler(req *http.Request) http.Handler {
	return n.findHandler(1, len(req.URL.Path)-1, req)
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
