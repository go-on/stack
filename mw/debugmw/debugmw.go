package debugmw

import (
	"fmt"
	"gopkg.in/go-on/stack.v1"
	"gopkg.in/go-on/stack.v1/mw"
	"gopkg.in/metakeule/backtrace.v1"
	"net/http"
)

type debugger struct {
	panicCodes []int
	*stack.Stack
	panicCatcher *backtrace.FmtPanicCatcher
}

func New(h http.Handler, skip, max int) *debugger {
	d := &debugger{}

	d.panicCatcher = backtrace.NewPanicCatcher(skip, max).SetFormat("    %s()\n    --> %s#%d\n\n")
	d.Stack = stack.New().
		Use(mw.Catch(d.catch)).
		Use(mw.Switch(d.decidePanicCode, d.pass, d.doPanic)).
		Use(mw.Prepare()).
		Use(mw.MethodOverride()).
		Use(mw.MethodOverrideByField("_method")).
		UseHandler(h)

	return d
}

func (d *debugger) SetPanicCodes(codes ...int) {
	d.panicCodes = codes
}

func (d *debugger) UnsetPanicCodes() {
	d.panicCodes = []int{}
}

func (d *debugger) decidePanicCode(r *http.Request) bool {
	return len(d.panicCodes) == 0
}

func (d *debugger) doPanic(w http.ResponseWriter, r *http.Request, next http.Handler) {
	mw.PanicCodes(d.panicCodes...).ServeHTTP(w, r, next)
}

func (d *debugger) pass(w http.ResponseWriter, r *http.Request, next http.Handler) {
	next.ServeHTTP(w, r)
}

func (d *debugger) catch(err interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ERROR on %s %s\n  Message: %v\n  Backtrace:\n", r.Method, r.URL.String(), err)
	d.panicCatcher.Catch(err, w, r)
}
