package stack

import (
	"net/http"
	"testing"
)

func mkRequestResponse() (w http.ResponseWriter, r *http.Request) {
	r, _ = http.NewRequest("GET", "/", nil)
	w = noHTTPWriter{}
	return
}

var wr = writeString("a")

func mkWrap(num int) http.Handler {
	var s Stack
	for i := 0; i < num; i++ {
		s.UseHandler(wr)
	}
	return s.Handler()
}

func times(n int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < n; i++ {
			wr.ServeHTTP(w, r)
		}
	})
}

func benchmark(h http.Handler, b *testing.B) {
	wr, req := mkRequestResponse()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(wr, req)
	}
}

func benchmarkWrapper(n int, b *testing.B) {
	b.StopTimer()
	h := mkWrap(n)
	benchmark(h, b)
}

func benchmarkSimple(n int, b *testing.B) {
	b.StopTimer()
	h := times(n)
	benchmark(h, b)
}

func BenchmarkServing2Simple(b *testing.B) {
	benchmarkSimple(2, b)
}

func BenchmarkServing2Wrappers(b *testing.B) {
	benchmarkWrapper(2, b)
}

func BenchmarkServing50Simple(b *testing.B) {
	benchmarkSimple(50, b)
}

func BenchmarkServing50Wrappers(b *testing.B) {
	benchmarkWrapper(50, b)
}

func BenchmarkServing100Simple(b *testing.B) {
	benchmarkSimple(100, b)
}

func BenchmarkServing100Wrappers(b *testing.B) {
	benchmarkWrapper(100, b)
}

func BenchmarkServing500Simple(b *testing.B) {
	benchmarkSimple(500, b)
}

func BenchmarkServing500Wrappers(b *testing.B) {
	benchmarkWrapper(500, b)
}

func BenchmarkServing1000Simple(b *testing.B) {
	benchmarkSimple(1000, b)
}

func BenchmarkServing1000Wrappers(b *testing.B) {
	benchmarkWrapper(1000, b)
}

func BenchmarkWrapping(b *testing.B) {
	b.StopTimer()
	var s Stack

	for i := 0; i < b.N; i++ {
		s.UseHandler(writeString("a"))
	}
	b.StartTimer()
	s.Handler()
}
