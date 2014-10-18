package middleware

import "net/http"

type stop struct{}

// ServeHTTP does nothing
func (stop) ServeHTTP(wr http.ResponseWriter, req *http.Request) {}

// Stop is a wrapper that does no processing but simply prevents further execution of
// next wrappers
var Stop = stop{}
