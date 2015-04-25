package stack

import "net/http"

// Middleware is like a http.Handler that is part of a chain of handlers handling the request
type Middleware interface {
	// ServeHTTP serves the request and may write to the given writer. Also may invoke the next handler
	ServeHTTP(writer http.ResponseWriter, request *http.Request, next http.Handler)
}

// ContextMiddleware is like Middleware that needs to be able to store and retrieve context bound to the request
type ContextMiddleware interface {
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler)
}

// ContextHandler is like a http.Handler that needs to be able to store and retrieve context bound to the request
type ContextHandler interface {
	ServeHTTP(ctx Contexter, wr http.ResponseWriter, req *http.Request)
}

// Wrapper wraps an http.Handler returning another http.Handler
type Wrapper interface {
	Wrap(http.Handler) http.Handler
}

// Swapper may swap between given value and the value inside a request context
type Swapper interface {
	// Swap must be defined on a pointer
	// and changes the the value of the pointer
	// to the value the replacement is pointing to
	Swap(replacement interface{})
}

// Contexter stores and retrieves per request data.
// Only one value per type can be stored.
type Contexter interface {

	// Set the given Swapper, replaces the value of the same type
	// Set may be run on the same Contexter concurrently
	Set(Swapper)

	// Get a given Swapper. If there is a Swapper of this type
	// inside the Contexter, the given Swappers Swap method is called with the stored value and true is returned.
	// If no Swapper of the same type could be found, false is returned
	// Get may be run on the same Contexter concurrently
	Get(Swapper) (has bool)

	// Del deletes a value of the given type.
	// Del may be run on the same Contexter concurrently
	Del(Swapper)

	// Transaction runs the given function inside a transaction. A TransactionContexter is passed to the
	// given function that might be used to call the Set, Get and Del methods inside the transaction.
	// However that methods must not be used concurrently
	Transaction(func(TransactionContexter))
}

// TransactionContexter stores and retrieves per request data via a hidden Contexter.
// It is meant to be used in the Transaction method of a Contexter.
// Only one TransactionContexter might be used at the same time for the same Contexter.
// No method of a TransactionContexter might be used concurrently
type TransactionContexter interface {
	// Set the given Swapper, replaces the value of the same type
	// Set may NOT be run on the same TransactionContexter concurrently
	Set(Swapper)

	// Get a given Swapper. If there is a Swapper of this type
	// inside the Contexter, the given Swappers Swap method is called with the stored value and true is returned.
	// If no Swapper of the same type could be found, false is returned
	// Get may NOT be run on the same TransactionContexter concurrently
	Get(Swapper) (has bool)

	// Del deletes a value of the given type.
	// Del may NOT be run on the same TransactionContexter concurrently
	Del(Swapper)
}
