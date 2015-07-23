/*
Package stacknosurf provides wrappers based on the github.com/justinas/nosurf package.

*/
package stacknosurf

import (
	"net/http"

	"gopkg.in/go-on/stack.v1"

	"gopkg.in/go-on/nosurf.v1"
)

// Token is the type that is saved inside a wrap.Contexter and
// represents a csrf token from the github.com/justinas/nosurf package.
type Token string

func (t *Token) Swap(repl interface{}) {
	*t = *(repl.(*Token))
}

// Tokenfield is the name of the form field that submits a csrf token
var TokenField = "csrf_token"

// SetToken is a middleware that sets a csrf token in the Contexter
// (response writer) on GET requests.
type SetToken struct{}

func (SetToken) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	if req.Method == "GET" {
		token := Token(nosurf.Token(req))
		ctx.Set(&token)
	}
	next.ServeHTTP(rw, req)
}

// CheckToken is a middleware that checks the token via the github.com/justinas/nosurf
// package. Its attributes relate to the corresponding nosurf options. If they are nil,
// they are not set.
type CheckToken struct {
	FailureHandler http.Handler
	BaseCookie     *http.Cookie
	ExemptPaths    []string
	ExemptGlobs    []string
	ExemptRegexps  []interface{}
	ExemptFunc     func(r *http.Request) bool
}

func (c *CheckToken) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	ns := nosurf.New(next)
	if c.BaseCookie != nil {
		ns.SetBaseCookie(*c.BaseCookie)
	}
	if c.FailureHandler != nil {
		ns.SetFailureHandler(c.FailureHandler)
	}

	if len(c.ExemptPaths) > 0 {
		ns.ExemptPaths(c.ExemptPaths...)
	}

	if len(c.ExemptGlobs) > 0 {
		ns.ExemptPaths(c.ExemptGlobs...)
	}

	if len(c.ExemptRegexps) > 0 {
		ns.ExemptRegexps(c.ExemptRegexps...)
	}

	if c.ExemptFunc != nil {
		ns.ExemptFunc(c.ExemptFunc)
	}
	ns.ServeHTTP(w, r)
}
