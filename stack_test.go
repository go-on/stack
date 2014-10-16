// +build go1.1

package stack

import (
	"fmt"
	"go/build"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestStack(t *testing.T) {
	tests := map[string]http.Handler{
		"abc": New(
			write("a"),
			write("b").ServeHTTPNext,
			write("c").ServeHTTP,
		),
		"ab": New(
			write("a").Wrap,
			writeStop("b"),
			write("c"),
		),
	}

	for body, h := range tests {
		rec, req := newTestRequest("GET", "/")
		h.ServeHTTP(rec, req)
		assertResponse(t, rec, body, 200)
	}
}

var gopath = filepath.SplitList(build.Default.GOPATH)[0]

func TestStackDebug(t *testing.T) {
	s := New(
		write("a"),
		writeDebug("b"),
		write("c").ServeHTTP,
	)

	file := filepath.Join(gopath, "src", "github.com", "go-on", "stack", "stack_test.go")
	if s.File != file {
		t.Errorf("s.File = %#v // expected %#v", s.File, file)
	}

	line := 43
	if s.Line != line {
		t.Errorf("s.Line = %#v // expected %#v", s.Line, line)
	}

	str := fmt.Sprintf("%s:%d", file, line)
	if !strings.Contains(s.String(), str) {
		t.Errorf("strings.Contains(%#v, %#v) == false", s.String(), str)
	}

	mwD := `stack.writeDebug = <writeDebug "b">`
	if got := s.Middleware[1]; got != mwD {
		t.Errorf("s.Middleware[1] == %#v != %#v", got, mwD)
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		in  interface{}
		out string
	}{
		{write("a"), "a"},
		{write("b").ServeHTTP, "b"},
		{write("c").Wrap, "c"},
		{writeStop("d"), "d"},
		{write("e").ServeHTTPNext, "e"},
		{http.HandlerFunc(write("f").ServeHTTP), "f"},
		{writeCtx, "ho"},
		{writeCtx2{}, "ho"},
	}
	_ = tests

	for _, test := range tests {
		h := Handler(test.in)
		rec := httptest.NewRecorder()
		New(
			Context,
			setCtx("ho"),
			h,
		).ServeHTTP(rec, nil)

		got := rec.Body.String()

		if got != test.out {
			t.Errorf("%#v != %#v", got, test.out)
		}
	}

}
