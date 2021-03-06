package stack

import (
	"net/http"
	"reflect"
	"sync"
)

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

func (c *context) Transaction(fn func(TransactionContexter)) {
	c.Lock()
	defer c.Unlock()
	fn(&contextTransaction{c})
}

type contextHandler struct {
	http.Handler
}

func (c *contextHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	ctx := &context{wr, map[interface{}]Swapper{}, &sync.RWMutex{}}
	ctx.Set(&ResponseWriter{wr})
	c.Handler.ServeHTTP(ctx, req)
}
