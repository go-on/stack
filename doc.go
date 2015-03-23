// Copyright (c) 2014 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package stack creates a fast and flexible middleware stack for http.Handlers.

Features

  - small and simple
  - based on http.Handler interface; nicely integrates with net/http
  - middleware stacks are http.Handlers too and may be embedded
  - has a solution for per request context sharing
  - easy debugging
  - low memory footprint
  - fast
  - no dependency apart from the standard library
  - easy to create adapters / wrappers for 3rd party middleware


Benchmarks (Go 1.3)

  // The overhead of n writes to http.ResponseWriter via n wrappers
  // vs n writes in a loop within a single http.Handler
  //
  BenchmarkServing2Simple    2000000         847 ns/op 1.00x
  BenchmarkServing2Wrappers  2000000         899 ns/op 1.06x

  BenchmarkServing100Simple    50000       40715 ns/op 1.00x
  BenchmarkServing100Wrappers  50000       42865 ns/op 1.05x

  BenchmarkServing50Simple    100000       20354 ns/op 1.00x
  BenchmarkServing50Wrappers  100000       21591 ns/op 1.06x


Accepted middleware

  // Functions
  func(next http.Handler) http.Handler
  func(http.ResponseWriter, *http.Request)
  func(http.ResponseWriter, *http.Request, next http.Handler)
  func(stack.Contexter, http.ResponseWriter, *http.Request)
  func(stack.Contexter, http.ResponseWriter, *http.Request, next http.Handler)

  // Interfaces
  ServeHTTP(http.ResponseWriter,*http.Request)                                      // http.Handler
  Wrap(http.Handler) http.Handler                                                   // stack.Wrapper
  ServeHTTP(http.ResponseWriter, *http.Request, next http.Handler)                  // stack.Middleware
  ServeHTTP(stack.Contexter, http.ResponseWriter, *http.Request)                    // stack.ContextHandler
  ServeHTTP(stack.Contexter, http.ResponseWriter, *http.Request, next http.Handler) // stack.ContextMiddleware

  // 3rd party middleware (via stack/adapter)
  Martini
  Negroni

  // 3rd party integrations
  abbot/go-http-auth (Basic and Digest Authentication)
  justinas/nosurf (CSRF protection)
  gorilla/sessions


Batteries included

Middleware can be found in the sub package stack/middleware.

Usage

  // define stack
  s := stack.New(
       // put your middleware here
       // from stacks point of view a http.Handler (may be the app or router) is just another middleware
  )

  // use it
  http.Handle("/", s)

Sharing Context

The ContextHandler method returns an http.Handler that wraps the given http.ResponseWriter with a stack.Contexter.

  // define stack
  s := stack.New(
       ...
  )

  // use it with a per request context
  http.Handle("/", s.ContextHandler())


To access the context inside a http.Handler or some middleware, the ResponseWriter must be type asserted to the stack.Contexter.

  func (rw http.ResponseWriter, req *http.Request) {
     ctx := rw.(stack.Contexter)
     ...
  }


Alternatively middlware might implement the stack.ContextHandler or stack.ContextMiddleware interface
or have the corresponding function signatures.

  func(stack.Contexter, http.ResponseWriter, *http.Request){}
  func(stack.Contexter, http.ResponseWriter, *http.Request, next http.Handler){}

  func (...) ServeHTTP(stack.Contexter, http.ResponseWriter, *http.Request){}
  func (...) ServeHTTP(stack.Contexter, http.ResponseWriter, *http.Request, next http.Handler)


To share per request context between middlewares, define a type for your data.
For each type only one value can be stored. The pointer of the type is stored and
must implement the Swapper interface.

  type MyVal string

  // Swap replaces the value of m with the value of val
  func (m *MyVal) Swap(val interface{}) {
     *m = *(val.(*MyVal)) // will never panic
  }

  type MyMiddleware struct{}

  func (mw *MyMiddleware) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
    var val MyVal = "some value"
    ctx.Set(&val)
    next.ServeHTTP(rw,req)
  }

  func writeVal(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
    var val MyVal

    if ctx.Get(&val) {
      fmt.Printf("value: %s\n", string(val))
    } else {
      fmt.Println("no value found")
    }
    next.ServeHTTP(rw, req)
  }

Don't forget to call the ContextHandler method on the outermost stack.
The returned handler must not be embedded inside another stack.

  func ExampleContextNew() {
    s := stack.New(
      stack.Context,
      writeVal,
      &MyMiddleware{},
      writeVal,
    )

    r, _ := http.NewRequest("GET", "/", nil)
    s.ContextHandler().ServeHTTP(nil, r)

    // Output:
    // no value found
    // value: some value
  }

Debugging

A Stack tracks the File and Line of its definition and saves middleware strings
that can have even more detailed information if the middleware is a fmt.Stringer.

see example_debug_test.go for an example

Original ResponseWriter

To get access to the original ResponseWriter there are several methods. Here is an example
of using the original ResponseWriter to type assert it to a http.Flusher.

  // access the original http.ResponseWriter to flush
  func flush(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request, next http.Handler) {
    var o stack.ResponseWriter
    ctx.Get(&o)

    // shortcut for the above
    // o := stack.ReclaimResponseWriter(rw)

    // type assertion of the original ResponseWriter
    if fl, is := o.ResponseWriter.(http.Flusher); is {
      fl.Flush()
    }

    // alternative shortcut for flushing, respecting an underlying ResponseWriter
    // stack.Flush(rw)
  }


  // shortcuts for underlying ResponseWriters

  Flush                 // flush a ResponseWriter if it is a http.Flusher
  CloseNotify           // returns a channel to notify if ResponseWriter is a http.CloseNotifier
  Hijack                // allows to hijack a connection if ResponseWriter is a http.Hijacker
  ReclaimResponseWriter // get the original ResponseWriter from a ResponseWriter with context

Other ResponseWriters

The package stack.responsewriter provides some ResponseWriter wrappers that help with development of middleware
and also support context sharing via embedding of a stack.Contexter if it is available.

FAQ

1. A ResponseWriter is an interface, because it may implement other interfaces from the http libary,
e.g. http.Flusher. If it is wrapped that underlying implementation is not accessible anymore

Answer: Since only one Contexter may be used within a stack, it is always possible to ask the
Contexter for the underlying ResponseWriter. This is what helper functions like
ReclaimResponseWriter(), Flush(), CloseNotify() and Hijack() do.

2. I get a panic "can't embed contextHandler" from the stack.New() function. Why?

Answer: You are trying to embed a stack that is wrapped by a contextHandler as middleware to another stack.
Simply remove the ContextHandler() method call from the inner stack as it is just allowed on the outermost stack.

Credits

Initial inspiration came from Christian Neukirchen's
rack for ruby some years ago.

Adapters come from carbocation/interpose (https://github.com/carbocation/interpose/blob/master/adaptors)
*/
package stack
