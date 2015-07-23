stack
=====
[ ![Codeship Status for go-on/stack](https://codeship.io/projects/7fa20300-38d9-0132-d3df-2a69fe4b0f90/status)](https://codeship.io/projects/42107)

[![Build status](https://ci.appveyor.com/api/projects/status/rv4pf8qwtj3n85vp?svg=true)](https://ci.appveyor.com/project/metakeule/stack)

Package stack creates a fast and flexible middleware stack for http.Handlers.

## Features

  - small and simple
  - based on http.Handler interface; nicely integrates with net/http
  - middleware stacks are http.Handlers too and may be embedded
  - has a solution for per request context sharing
  - easy debugging
  - low memory footprint
  - fast
  - no dependency apart from the standard library
  - easy to create adapters / wrappers for 3rd party middleware

## Status

Not ready for general consumption yet. API may change at any time.

## Benchmarks (Go 1.4)

The overhead of n writes to http.ResponseWriter via n wrappers vs n writes in a loop within a single http.Handler on my laptop

```
  BenchmarkServing2Simple       10000000     201 ns/op 1.00x 
  BenchmarkServing2Wrappers      5000000     250 ns/op 1.24x
  
  BenchmarkServing50Simple        300000    4833 ns/op 1.00x
  BenchmarkServing50Wrappers      200000    6503 ns/op 1.35x
  
  BenchmarkServing100Simple       200000    9810 ns/op 1.00x
  BenchmarkServing100Wrappers     100000   13489 ns/op 1.38x
  
  BenchmarkServing500Simple        30000   48970 ns/op 1.00x
  BenchmarkServing500Wrappers      20000   75562 ns/op 1.54x
  
  BenchmarkServing1000Simple       20000   98737 ns/op 1.00x
  BenchmarkServing1000Wrappers     10000  163671 ns/op 1.66x
```

## Accepted middleware

```
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
```

## Batteries included

Middleware can be found in the sub package stack/mw.

## Router

A simple router is included in the sub package router.

## Usage

```go
  // define middleware
  func middleware(w http.ResponseWriter, r *http.Request, next http.Handler) {
    w.Write([]byte("middleware speaking\n"))
    next.ServeHTTP(w,r)
  }

  // define app
  func app(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("app here"))
  }  

  // define stack
  var s Stack
  s.UseFunc(middleware)
  // add other middleware for each supported type there is a UseXXX method
  
  // use it
  http.Handle("/", s.WrapFunc(app))
```

## Sharing Context

Context is shared by wrapping the http.ResponseWriter with another one that also implements the Contexter interface. This new ResponseWriter is then passed to the middleware. In order to use the context the middleware must have
the Contexter as first parameter. 
Context data can be any type that has a pointer method implementing the Recoverer interface. Each type can only be saved once per request.

```go
  // define context data type, may be any type
  type Name string

  // implement Recoverer interface
  // Recover replaces the value of m with the value of val
  func (n *Name) Recover(val interface{}) {
    *n = *(val.(*Name)) // will never panic
  }

  // define middleware that uses context
  func middleware(ctx stack.Contexter, w http.ResponseWriter, r *http.Request, next http.Handler) {
    var n Name = "Hulk"
    ctx.Set(&n)
    next.ServeHTTP(w, r)
  }

  // define app that uses context
  func app(ctx stack.Contexter, w http.ResponseWriter, r *http.Request) {
    var n Name
    ctx.Get(&n)
    w.Write([]byte("name is " + string(n)))
  }

  // define stack
  s := stack.New()
  s.UseFuncWithContext(middleware)
  // add other middleware for each supported type there is a UseXXXWithContext method
       

  // make sure to call WrapWithContext or WrapFuncWithContext to create the context object
  http.Handle("/", s.WrapFuncWithContext(app))
```

## Original ResponseWriter

To get access to the original ResponseWriter there are several methods. Here is an example of using the original ResponseWriter to type assert it to a http.Flusher.

```go
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
```
```
  // shortcuts for underlying ResponseWriters

  Flush                 // flush a ResponseWriter if it is a http.Flusher
  CloseNotify           // returns a channel to notify if ResponseWriter is a http.CloseNotifier
  Hijack                // allows to hijack a connection if ResponseWriter is a http.Hijacker
  ReclaimResponseWriter // get the original ResponseWriter from a ResponseWriter with context
```

## Other ResponseWriters

The package stack.responsewriter provides some ResponseWriter wrappers that help with development of middleware
and also support context sharing via embedding of a stack.Contexter if it is available.

## Server

A simple ready-to-go server is inside the server subpackage.

```go

package main

import (
  "github.com/go-on/stack/server"
  "github.com/go-on/stack/mw"
)

func main() {
  s := server.New()
  s.INDEX("hu", mw.TextString("hu"))
  s.ListenAndServe(":8080")
}

```

## Credits

Initial inspiration came from Christian Neukirchen's
rack for ruby some years ago.

Adapters come from [carbocation/interpose](https://github.com/carbocation/interpose/blob/master/adaptors)