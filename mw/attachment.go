package mw

import (
	// "fmt"
	"net/http"
)

type Attachment string

func (a Attachment) Set(w http.ResponseWriter) {
	val := "attachment"
	if string(a) != "" {
		val += `; filename=` + string(a)
	}
	// fmt.Println("setting Content-Disposition to ", val)
	w.Header().Set("Content-Disposition", val)
}

func (a Attachment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Set(w)
}
