package main

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func BenchmarkRouter(b *testing.B) {
	router := httprouter.New()
	router.GET("/v1", docs)
	router.GET("/v1/users/{id}", users)

	req, err := http.NewRequest("GET", "http://localhost:8080/v1/users/5", nil)
	if err != nil {
		b.Fatal(err)
	}

	w := &mockResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
