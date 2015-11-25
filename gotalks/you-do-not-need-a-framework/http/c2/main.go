package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func main() {
	http.Handle("/hello", alice.New(Timing(), Method("GET")).ThenFunc(hello))
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Timing() alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			started := time.Now()
			h.ServeHTTP(w, req)
			log.Printf("%s %s in %s\n", req.Method, req.URL.String(), time.Since(started))
		})
	}
}
func Method(method string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method != method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, http.StatusText(http.StatusMethodNotAllowed))
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}
