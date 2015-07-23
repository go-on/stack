/*
Package rest provides a simple and opinionated rest rest

Principles

1. RESTs must not know about the route they are mounted to. Mount rest from the outside just like handlers.
2. Leave dispatching based on paths to extra muxers, like http.ServeMux
3. Route paths may only have one parameter or no parameter
4. If a route paths ends with a slash, there is no parameter.
5. If a route path does not end with a slash, the parameter is the trimmed part of the url path from the last slash to the end
6. Other parameters must be part of the url query params or be inside the request body.
7. There is just one muxer that dispatches based on the url path and it is used at the top level
8. There are no sub rests / no hierarchie of rests. Each rest is independent from the other.
9. The mounting of handlers on paths and subpaths must be defined inside an app.
10. There is an example how that could be done.
11. Each rest can only have a max of 7 handlers: One for the methods GET, PUT, PATCH, DELETE, HEAD and POST and an index handler
12. The index handler is called for GET requests with a path ending with a slash.

This leads to

1. Decoupled rests that can easily be mounted anywhere.
2. Fast and flexible mount paths without a need for placeholders, parsing or regular expressions.
3. Parameter type conversion inside the handlers where it belongs to.
4. Flexible integration of middleware around and within rests.
5. URLs are simply strings that could be stored in variables and parametrized by concatenation.
6. No special knowledge about url handling and rest internals required.

*/
package rest

import (
	"net/http"
	"strings"
)

const (
	index = iota
	get
	post
	put
	del
	patch
	head
)

// REST is an array of maximal 6 handlers
// One for the methods GET, PUT, PATCH, DELETE, HEAD and POST and the index handler
// Each handler may be missing; then a 404 status will be written for requests of that method
type REST [7]http.Handler

// Get registers a handler for the GET method and an url path that has a parameter (does not end with a slash)
func (r REST) Get(h http.Handler) REST { r[get] = h; return r }

// Head registers a handler for the HEAD method and an url path that has a parameter (does not end with a slash)
func (r REST) Head(h http.Handler) REST { r[head] = h; return r }

// Index registers a handler for the GET method and an url path that ends with a slash (no parameter)
func (r REST) Index(h http.Handler) REST { r[index] = h; return r }

// Post registers a handler for the POST method and an url path that ends with a slash (no parameter)
func (r REST) Post(h http.Handler) REST { r[post] = h; return r }

// Put registers a handler for the PUT method and an url path that has a parameter (does not end with a slash)
func (r REST) Put(h http.Handler) REST { r[put] = h; return r }

// Patch registers a handler for the PATCH method and an url path that has a parameter (does not end with a slash)
func (r REST) Patch(h http.Handler) REST { r[patch] = h; return r }

// Delete registers a handler for the DELETE method and an url path that has a parameter (does not end with a slash)
func (r REST) Delete(h http.Handler) REST { r[del] = h; return r }

// GetFunc is like Get but for http.HandlerFunc
func (r REST) GetFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Get(http.HandlerFunc(fn))
}

// HeadFunc is like Head but for http.HandlerFunc
func (r REST) HeadFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Head(http.HandlerFunc(fn))
}

// IndexFunc is like Index but for http.HandlerFunc
func (r REST) IndexFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Index(http.HandlerFunc(fn))
}

// PostFunc is like Post but for http.HandlerFunc
func (r REST) PostFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Post(http.HandlerFunc(fn))
}

// PutFunc is like Put but for http.HandlerFunc
func (r REST) PutFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Put(http.HandlerFunc(fn))
}

// PatchFunc is like Patch but for http.HandlerFunc
func (r REST) PatchFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Patch(http.HandlerFunc(fn))
}

// DeleteFunc is like Delete but for http.HandlerFunc
func (r REST) DeleteFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	return r.Delete(http.HandlerFunc(fn))
}

// Get creates a REST and registers the handler via Get method call
func Get(h http.Handler) REST { var r REST; return r.Get(h) }

// Head creates a REST and registers the handler via Head method call
func Head(h http.Handler) REST { var r REST; return r.Head(h) }

// Index creates a REST and registers the handler via Index method call
func Index(h http.Handler) REST { var r REST; return r.Index(h) }

// Post creates a REST and registers the handler via Post method call
func Post(h http.Handler) REST { var r REST; return r.Post(h) }

// Put creates a REST and registers the handler via Put method call
func Put(h http.Handler) REST { var r REST; return r.Put(h) }

// Patch creates a REST and registers the handler via Patch method call
func Patch(h http.Handler) REST { var r REST; return r.Patch(h) }

// Delete creates a REST and registers the handler via Delete method call
func Delete(h http.Handler) REST { var r REST; return r.Delete(h) }

// GetFunc is like Get but for http.HandlerFunc
func GetFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.GetFunc(fn)
}

// HeadFunc is like Head but for http.HandlerFunc
func HeadFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.HeadFunc(fn)
}

// IndexFunc is like Index but for http.HandlerFunc
func IndexFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.IndexFunc(fn)
}

// PostFunc is like Post but for http.HandlerFunc
func PostFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.PostFunc(fn)
}

// PutFunc is like Put but for http.HandlerFunc
func PutFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.PutFunc(fn)
}

// PatchFunc is like Patch but for http.HandlerFunc
func PatchFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.PatchFunc(fn)
}

// DeleteFunc is like Delete but for http.HandlerFunc
func DeleteFunc(fn func(wr http.ResponseWriter, req *http.Request)) REST {
	var r REST
	return r.DeleteFunc(fn)
}

func (r REST) has(method int) bool {
	return r[method] != nil
}

// Options returns the options (method) for index or non index requests
func (r REST) Options(isIndex bool) (o []string) {
	o = append(o, "OPTIONS")

	if isIndex {
		if r.has(index) {
			o = append(o, "GET")
		}
		if r.has(post) {
			o = append(o, "POST")
		}
		return
	}

	if r.has(get) {
		o = append(o, "GET")
	}
	if r.has(patch) {
		o = append(o, "PATCH")
	}

	if r.has(put) {
		o = append(o, "PUT")
	}

	if r.has(head) {
		o = append(o, "HEAD")
	}

	if r.has(patch) {
		o = append(o, "PATCH")
	}

	if r.has(del) {
		o = append(o, "DELETE")
	}

	return
}

// IsIndex checks if the given request is an index request
func IsIndex(req *http.Request) bool {
	s := req.URL.Path
	return s[len(s)-1] == '/'
}

// ServeOptions serves the OPTIONS request
func (r REST) ServeOptions(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Allow", strings.Join(r.Options(IsIndex(req)), ","))
	wr.WriteHeader(http.StatusOK)
}

// Handler returns the handler for the given method and index or non index.
// found is true if a handler could be found
func (r REST) Handler(method string, isIndex bool) (h http.Handler, found bool) {
	if method == "OPTIONS" {
		return http.HandlerFunc(r.ServeOptions), true
	}

	if isIndex {
		switch method {
		case "GET":
			h = r[index]
		case "POST":
			h = r[post]
		}
		found = h != nil
		return
	}
	switch method {
	case "GET":
		h = r[get]
	case "PATCH":
		h = r[patch]
	case "PUT":
		h = r[put]
	case "DELETE":
		h = r[del]
	case "HEAD":
		h = r[head]
	}
	found = h != nil
	return
}

func (r REST) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if h, found := r.Handler(req.Method, IsIndex(req)); found {
		h.ServeHTTP(wr, req)
		return
	}
	http.NotFound(wr, req)
}

func Param(req *http.Request) string {
	if len(req.URL.Path) < 2 {
		return ""
	}
	return strings.TrimSpace(req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:])
}
