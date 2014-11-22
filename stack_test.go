// +build go1.1

package stack

import (
	"fmt"
	"go/build"
	"net/http"
	"path/filepath"
	"testing"
)

func init() {
	fmt.Println(build.Default.BuildTags)
}

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
