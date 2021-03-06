/*
The MIT License (MIT)

Copyright (c) 2014 carbocation

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
Package adapter_negroni is an adapter for negroni middleware to be used with github.com/go-on/stack.
It is a copy of https://github.com/carbocation/interpose/blob/master/adaptors/negroni.go
*/
package adapter_negroni

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

func FromNegroni(handler negroni.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		n := negroni.New()
		n.Use(handler)
		n.UseHandler(next)
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			n.ServeHTTP(rw, req)
		})
	}
}

func HandlerFromNegroni(handler negroni.Handler) http.Handler {
	n := negroni.New()
	n.Use(handler)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		n.ServeHTTP(rw, req)
	})
}
