package responsewriter

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-on/stack"
)

// JSONDecoder is a ResponseWriter wrapper that decodes everything writen to it to a json object.
//
type JSONDecoder struct {
	// ResponseWriter is the underlying response writer that is wrapped by Buffer
	http.ResponseWriter

	// if the underlying ResponseWriter is a Contexter, that Contexter is saved here
	stack.Contexter

	// Buffer is the underlying io.Writer that buffers the response body
	bf bytes.Buffer

	// Code is the cached status code
	Code int

	// changed tracks if anything has been set on the responsewriter. Also reads from the header
	// are seen as changes
	changed bool

	// header is the cached header
	header http.Header
}

// make sure to fulfill the Contexter interface
var _ stack.Contexter = &JSONDecoder{}

// NewJSONDecoder creates a new JSONDecoder by wrapping the given response writer.
func NewJSONDecoder(w http.ResponseWriter) (dec *JSONDecoder) {
	dec = &JSONDecoder{}
	dec.ResponseWriter = w
	if ctx, ok := dec.ResponseWriter.(stack.Contexter); ok {
		dec.Contexter = ctx
	} else {
		panic("no contexter")
	}
	dec.header = make(http.Header)
	return
}

// Header returns the cached http.Header and tracks this call as change
func (d *JSONDecoder) Header() http.Header {
	d.changed = true
	return d.header
}

// WriteHeader writes the cached status code and tracks this call as change
func (d *JSONDecoder) WriteHeader(i int) {
	d.changed = true
	d.Code = i
}

// Write writes to the underlying buffer and tracks this call as change
func (d *JSONDecoder) Write(b []byte) (int, error) {
	d.changed = true
	return d.bf.Write(b)
}

// Decode decodes the buffer into the given target
func (d *JSONDecoder) Decode(target interface{}) error {
	return json.NewDecoder(&d.bf).Decode(target)
}

// FlushCode flushes the status code to the underlying responsewriter if it was set.
func (d *JSONDecoder) FlushCode() {
	if d.Code != 0 {
		d.ResponseWriter.WriteHeader(d.Code)
	}
}

// HasChanged returns true if Header, WriteHeader or Write has been called
func (d *JSONDecoder) HasChanged() bool {
	return d.changed
}

// FlushHeaders adds the headers to the underlying ResponseWriter, removing them from Buffer.
func (d *JSONDecoder) FlushHeaders() {
	header := d.ResponseWriter.Header()
	for k, v := range d.header {
		header.Del(k)
		for _, val := range v {
			header.Add(k, val)
		}
	}
}
