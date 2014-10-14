package stack

// Contexter stores and retrieves data based on their types.

import (
	"net/http"
	"reflect"
	"sync"
)

// Only one value per type can be stored inside a contexter
type Contexter interface {

	// Set the given Swapper, replaces the value of the same type
	// Set may be run on the same Contexter concurrently
	Set(Swapper)

	// Get a given Swapper. If there is a Swapper of this type
	// inside the Contexter, the given Swappers Swap called and given the stored Swapper and true is returned.
	// If no Swapper of the same type could be found, false is returned
	// Get may be run on the same Contexter concurrently
	Get(Swapper) (has bool)

	// Del deletes a value of the given type
	// Del may be run on the same Contexter concurrently
	Del(Swapper)

	// Transaction runs the given function inside a transaction. A new contexter is passed to the
	// given function that might be used to call the Set, Get and Del methods inside the transaction.
	// However that Contexers methods must not be used concurrently
	Transaction(func(Contexter))
}

type contextTransaction struct {
	*context
}

func (c *contextTransaction) Set(val Swapper) {
	c.data[reflect.TypeOf(val)] = val
}

func (c *contextTransaction) Del(val Swapper) {
	delete(c.data, reflect.TypeOf(val))
}

func (c *contextTransaction) Get(target Swapper) bool {
	src, has := c.data[reflect.TypeOf(target)]
	if !has {
		return false
	}
	target.Swap(src)
	return true
}

var _ Contexter = &contextTransaction{}
var _ Contexter = &context{}

func (c *contextTransaction) Transaction(fn func(Contexter)) {
	panic("nested transactions are not allowed")
}

type context struct {
	http.ResponseWriter // you always need this
	data                map[interface{}]Swapper
	*sync.RWMutex
}

func (c *context) Set(val Swapper) {
	c.Lock()
	defer c.Unlock()
	c.data[reflect.TypeOf(val)] = val
}

func (c *context) Del(val Swapper) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, reflect.TypeOf(val))
}

func (c *context) Get(target Swapper) bool {
	c.RLock()
	defer c.RUnlock()
	src, has := c.data[reflect.TypeOf(target)]
	if !has {
		return false
	}
	target.Swap(src)
	return true
}

func (c *context) Transaction(fn func(Contexter)) {
	c.Lock()
	defer c.Unlock()
	fn(&contextTransaction{c})
}

// Context calls the next http.Handlers ServeHTTP method
// with a ResponseWriter that implements the Contexter interface.
// If the ResponseWriter already implements the Contexter interface, it is simply passed through.
func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		_, has := wr.(Contexter)
		if has {
			next.ServeHTTP(wr, req)
			return
		}
		ctx := &context{wr, map[interface{}]Swapper{}, &sync.RWMutex{}}
		ctx.Set(&ResponseWriter{wr})
		next.ServeHTTP(ctx, req)
	})
}
