/*
Package router provides a simple router that nicely integrates with stacks and middleware.

*/
package router

import (
	"github.com/go-on/stack"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Router is a mountable http.Handler that dispatches requests to other http.Handlers
// by request.Method and path segment.
type Router struct {
	prefix       string
	parent       *Router
	getRoutes    map[string]http.Handler
	onceGet      sync.Once
	postRoutes   map[string]http.Handler
	oncePost     sync.Once
	deleteRoutes map[string]http.Handler
	onceDelete   sync.Once
	patchRoutes  map[string]http.Handler
	oncePatch    sync.Once
	putRoutes    map[string]http.Handler
	oncePut      sync.Once
	subHandler   map[string]http.Handler
	onceSub      sync.Once
	notFound     http.Handler
}

// New returns a new router
func New() *Router {
	return &Router{}
}

func (r *Router) SetNotFound(h http.Handler) *Router {
	r.notFound = h
	return r
}

// MountPath returns a path that is prefixed to each path segment of a route.
func (r *Router) MountPath() string {
	if r.parent != nil {
		return path.Join(r.parent.MountPath(), r.prefix)
	}
	if r.prefix == "" {
		return "/"
	}
	return r.prefix
}

// SegmentMatch is the regular expression to which each path segment must match
var SegmentMatch = regexp.MustCompile("^[a-z][-.a-z0-9]*$")

// validateSegment returns an error message, if the pathSegment is invalid and an empty string if valid
func validateSegment(pathSegment string) (s string) {
	if !SegmentMatch.MatchString(pathSegment) {
		s = "invalid pathSegment syntax, must match ^[a-z][-.a-z0-9]*$"
	}
	return
}

// registerSubHandler registers the given http.Handler for the path segment regardless of the http method
func (r *Router) registerSubHandler(pathSegment string, h http.Handler) {
	if msg := validateSegment(pathSegment); msg != "" {
		panic(msg)
	}

	if _, has := r.subHandler[pathSegment]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}
	if _, has := r.getRoutes[pathSegment]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}
	if _, has := r.getRoutes[pathSegment+"/:"]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}
	if _, has := r.postRoutes[pathSegment]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}
	if _, has := r.putRoutes[pathSegment+"/:"]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}

	if _, has := r.patchRoutes[pathSegment+"/:"]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}
	if _, has := r.deleteRoutes[pathSegment+"/:"]; has {
		panic("route for pathSegment " + pathSegment + " already defined")
	}

	r.onceSub.Do(func() { r.subHandler = map[string]http.Handler{} })
	r.subHandler[pathSegment] = h
}

// Handle registers a http.Handler for the path segment regardless of the http method
func (r *Router) Handle(pathSegment string, h http.Handler) {
	if _, isRouter := h.(*Router); isRouter {
		panic("handler must not be a *Router, use Mount() method to mount a router")
	}

	r.registerSubHandler(pathSegment, h)
}

// HandleFunc is a shortcut for Handle(pathSegment, http.HandlerFunc(fn))
func (r *Router) HandleFunc(pathSegment string, h http.HandlerFunc) {
	r.Handle(pathSegment, h)
}

// Mount mounts the router beneath the given parent under the given pathSegment.
// A router may only be mounted once. But a router may be parent of many subrouters as long as
// they have different pathSegmentes. The pathSegment must match to SegmentMatch.
// Any violation of this rules results in a panic and mounting is meant to be done on startup.
// The optional middleware will surround the router
func (r *Router) Mount(pathSegment string, parent *Router, middleware ...interface{}) {
	if r.prefix != "" {
		panic("cannot mount the same sub router twice")
	}

	var h http.Handler

	if len(middleware) > 0 {
		middleware = append(middleware, r)
		h = stack.New(middleware...)
	} else {
		h = r
	}
	parent.registerSubHandler(pathSegment, h)

	r.parent = parent
	r.prefix = pathSegment
}

// route validates the pathSegment, registers the handler and returns a function that returns the URL
func (r *Router) route(method string, pathSegment string, h http.Handler) func() string {
	if msg := validateSegment(pathSegment); msg != "" {
		panic(msg)
	}

	if len(r.subHandler) > 0 {
		if _, has := r.subHandler[pathSegment]; has {
			panic("route for pathSegment already defined")
		}
	}
	r.registerHandler(method, pathSegment, h)
	return func() string {
		return path.Join(r.MountPath(), pathSegment)
	}
}

// routeParam validates the pathSegment, registers the handler and returns a function that returns the URL for a given param
func (r *Router) routeParam(method string, pathSegment string, h http.Handler) func(string) string {
	if msg := validateSegment(pathSegment); msg != "" {
		panic(msg)
	}

	if len(r.subHandler) > 0 {
		if _, has := r.subHandler[pathSegment]; has {
			panic("route for pathSegment already defined")
		}
	}
	r.registerHandler(method, pathSegment+"/:", h)
	return func(param string) string {
		return path.Join(r.MountPath(), pathSegment, param)
	}
}

// registerHandler registers the given handler for the given prefix and the given http method
func (r *Router) registerHandler(method string, prefix string, h http.Handler) {
	var mp map[string]http.Handler

	switch method {
	case "GET":
		r.onceGet.Do(func() { r.getRoutes = map[string]http.Handler{} })
		mp = r.getRoutes
	case "POST":
		r.oncePost.Do(func() { r.postRoutes = map[string]http.Handler{} })
		mp = r.postRoutes
	case "PATCH":
		r.oncePatch.Do(func() { r.patchRoutes = map[string]http.Handler{} })
		mp = r.patchRoutes
	case "PUT":
		r.oncePut.Do(func() { r.putRoutes = map[string]http.Handler{} })
		mp = r.putRoutes
	case "DELETE":
		r.onceDelete.Do(func() { r.deleteRoutes = map[string]http.Handler{} })
		mp = r.deleteRoutes
	default:
		panic("unknown method " + method)
	}

	if _, has := mp[prefix]; has {
		panic("route for pathSegment already defined")
	}

	mp[prefix] = h
}

// ServeHTTP serves the http request by delegating to the proper handler for the given
// http method and path
func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	prefix := rq.URL.Path
	prefix = strings.TrimLeft(prefix, "/")
	prefix = strings.TrimRight(prefix, "/")
	r.serveRoute(prefix, w, rq)
}

// findHandler returns the handler for the given prefix and the correct request method
// if there is a parameter, it will be saved inside the request.URL.Fragment
func (r *Router) findHandler(rq *http.Request, prefix string) http.Handler {
	var m map[string]http.Handler
	switch rq.Method {
	case "GET":
		m = r.getRoutes
	case "POST":
		m = r.postRoutes
	case "PUT":
		m = r.putRoutes
	case "PATCH":
		m = r.patchRoutes
	case "DELETE":
		m = r.deleteRoutes
	}

	if handler, found := m[prefix]; found {
		return handler
	}

	pos := strings.LastIndex(prefix, "/")
	if pos == -1 || pos+1 >= len(prefix) {
		return nil
	}

	param := string(prefix[pos+1:])
	rq.URL.Fragment = param
	prefix = string(prefix[:pos] + "/:")

	return m[prefix]
}

// serve the route for a request
func (r *Router) serveRoute(prefix string, w http.ResponseWriter, rq *http.Request) {
	var pr = prefix
	if idx := strings.Index(prefix, "/"); idx > 0 {
		pr = prefix[:idx]
	}

	if sub, found := r.subHandler[pr]; found {
		rq.URL.Path = strings.TrimPrefix(rq.URL.Path, "/"+pr)
		sub.ServeHTTP(w, rq)
		return
	}

	var handler = r.findHandler(rq, prefix)

	if handler == nil {
		if r.notFound != nil {
			r.notFound.ServeHTTP(w, rq)
			return
		}

		if rq.Method == "GET" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`not found`))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`method not allowed`))
		}
		return
	}
	handler.ServeHTTP(w, rq)
	return
}

// GET registers a http.Handler for the given path segment plus an id and the GET http method. It returns a function that returns the
// URL for a given item id. When the request is served, the item id will be saved inside the request.URL.Fragment
// and may be used by the http.Handler
func (r *Router) GET(pathSegment string, h http.Handler) func(id string) string {
	return r.routeParam("GET", pathSegment, h)
}

// POST registers a http.Handler for the given path segment plus an id and the POST http method. It returns a function that returns the
// URL.
func (r *Router) POST(pathSegment string, h http.Handler) func() string {
	return r.route("POST", pathSegment, h)
}

// PUT registers a http.Handler for the given path segment plus an id and the PUT http method. It returns a function that returns the
// URL for a given item id. When the request is served, the item id will be saved inside the request.URL.Fragment
// and may be used by the http.Handler
func (r *Router) PUT(pathSegment string, h http.Handler) func(id string) string {
	return r.routeParam("PUT", pathSegment, h)
}

// PATCH registers a http.Handler for the given path segment plus an id and the PATCH http method. It returns a function that returns the
// URL for a given item id. When the request is served, the item id will be saved inside the request.URL.Fragment
// and may be used by the http.Handler
func (r *Router) PATCH(pathSegment string, h http.Handler) func(id string) string {
	return r.routeParam("PATCH", pathSegment, h)
}

// DELETE registers a http.Handler for the given path segment plus an id and the DELETE http method. It returns a function that returns the
// URL for a given item id. When the request is served, the item id will be saved inside the request.URL.Fragment
// and may be used by the http.Handler
func (r *Router) DELETE(pathSegment string, h http.Handler) func(id string) string {
	return r.routeParam("DELETE", pathSegment, h)
}

// INDEX registers a http.Handler for the given path segment and the GET http method. It returns a function that returns the
// URL. Note that here is the exception that the path segment may be the empty string. Then the handler will be the main handler
// for the router.
func (r *Router) INDEX(pathSegment string, h http.Handler) func() string {
	if pathSegment != "" {
		if msg := validateSegment(pathSegment); msg != "" {
			panic(msg)
		}
	}

	if len(r.subHandler) > 0 {
		if _, has := r.subHandler[pathSegment]; has {
			panic("route for pathSegment already defined")
		}
	}
	r.registerHandler("GET", pathSegment, h)
	return func() string {
		return path.Join(r.MountPath(), pathSegment)
	}
}

// GETFunc is a shortcut for GET(pathSegment, http.HandlerFunc(fn))
func (r *Router) GETFunc(pathSegment string, h http.HandlerFunc) func(id string) string {
	return r.GET(pathSegment, h)
}

// POSTFunc is a shortcut for POST(pathSegment, http.HandlerFunc(fn))
func (r *Router) POSTFunc(pathSegment string, h http.HandlerFunc) func() string {
	return r.POST(pathSegment, h)
}

// PUTFunc is a shortcut for PUT(pathSegment, http.HandlerFunc(fn))
func (r *Router) PUTFunc(pathSegment string, h http.HandlerFunc) func(id string) string {
	return r.PUT(pathSegment, h)
}

// PATCHFunc is a shortcut for PATCH(pathSegment, http.HandlerFunc(fn))
func (r *Router) PATCHFunc(pathSegment string, h http.HandlerFunc) func(id string) string {
	return r.PATCH(pathSegment, h)
}

// DELETEFunc is a shortcut for DELETE(pathSegment, http.HandlerFunc(fn))
func (r *Router) DELETEFunc(pathSegment string, h http.HandlerFunc) func(id string) string {
	return r.DELETE(pathSegment, h)
}

// INDEXFunc is a shortcut for INDEX(pathSegment, http.HandlerFunc(fn))
func (r *Router) INDEXFunc(pathSegment string, h http.HandlerFunc) func() string {
	return r.INDEX(pathSegment, h)
}

type Getter interface {
	GET(http.ResponseWriter, *http.Request)
}

type Indexer interface {
	INDEX(http.ResponseWriter, *http.Request)
}

type Poster interface {
	POST(http.ResponseWriter, *http.Request)
}

type Patcher interface {
	PATCH(http.ResponseWriter, *http.Request)
}

type Putter interface {
	PUT(http.ResponseWriter, *http.Request)
}

type Deleter interface {
	DELETE(http.ResponseWriter, *http.Request)
}

func (r *Router) RouteObject(object interface{}, pathSegment string) (url func() string, urlParam func(string) string) {
	if o, ok := object.(Getter); ok {
		r.GETFunc(pathSegment, o.GET)
	}

	if o, ok := object.(Indexer); ok {
		r.INDEXFunc(pathSegment, o.INDEX)
	}

	if o, ok := object.(Poster); ok {
		r.POSTFunc(pathSegment, o.POST)
	}

	if o, ok := object.(Putter); ok {
		r.PUTFunc(pathSegment, o.PUT)
	}

	if o, ok := object.(Patcher); ok {
		r.PATCHFunc(pathSegment, o.PATCH)
	}

	if o, ok := object.(Deleter); ok {
		r.DELETEFunc(pathSegment, o.DELETE)
	}

	fn1 := func() string { return path.Join(r.MountPath(), pathSegment) }
	fn2 := func(param string) string { return path.Join(r.MountPath(), pathSegment, param) }
	return fn1, fn2
}

// RelativePath returns the path relative to the mount path (makes only sense for path beneath the mountpath)
func (r *Router) RelativePath(path string) (rel string, err error) {
	return filepath.Rel(r.MountPath(), path)
}
