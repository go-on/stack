package middleware

// JSONString is an utf-8 string that is a http.Handler

import (
	"encoding/json"

	"github.com/go-on/stack"

	"net/http"
)

func WriteJSON(ctx stack.Contexter, rw http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)

	if err != nil {
		SetError(err, ctx)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(b)
}

type json_ struct {
	obj interface{}
}

func (j *json_) ServeHTTP(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request) {
	WriteJSON(ctx, rw, j.obj)
}

func JSON(obj interface{}) *json_ {
	return &json_{obj}
}
