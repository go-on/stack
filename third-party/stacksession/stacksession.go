/*
Package stacksession provides wrappers based on the github.com/gorilla/sessions.
It assumes one session per request. If you need multiple sessions pre request, define your own
wrapper or use gorilla/sessions directly.

stacksession expects the ResponseWriter to a wrap.Contexter supporting sessions.Session and error

*/
package stacksession

import (
	"net/http"

	"gopkg.in/go-on/stack.v6/mw"

	"gopkg.in/go-on/sessions.v1"
	"gopkg.in/go-on/stack.v6"
)

type store struct {
	Store sessions.Store
	Name  string
}

type Session struct {
	sessions.Session
}

func (s *Session) Reclaim(repl interface{}) {
	*s = *(repl.(*Session))
}

func (s *store) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	sess, _ := s.Store.Get(req, s.Name)
	if sess != nil {
		ss := Session{*sess}
		ctx.Set(&ss)
	}
	next.ServeHTTP(rw, req)
}

func NewStore(st sessions.Store, name string) *store {
	return &store{Store: st, Name: name}
}

type saveAndClear struct{}

func (s saveAndClear) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
	defer func() {
		var sess Session
		if ctx.Get(&sess) {
			err := sess.Save(req, rw)
			if err != nil {
				mw.SetError(err, ctx)
			}
		}
	}()
	next.ServeHTTP(rw, req)
}

// SaveAndClear saves the session and clears up any references to the request.
// it should be used at the beginning of the mw chain (but after the Contexter)
var SaveAndClear = saveAndClear{}
