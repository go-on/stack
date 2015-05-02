package mw

import (
	"net/http"
)

type switcher struct {
	mw     [2]func(w http.ResponseWriter, r *http.Request, next http.Handler)
	decide func(r *http.Request) bool
}

func (s *switcher) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if s.decide(r) {
		s.mw[0](w, r, next)
		return
	}
	s.mw[1](w, r, next)
}

// Switch switches between 2 middleware functions via call of a decider function
// If the decider function returns true, truemw is run, otherwise the falsemw
func Switch(decider func(r *http.Request) bool, truemw, falsemw func(w http.ResponseWriter, r *http.Request, next http.Handler)) *switcher {
	return &switcher{
		mw:     [2]func(w http.ResponseWriter, r *http.Request, next http.Handler){truemw, falsemw},
		decide: decider,
	}
}
