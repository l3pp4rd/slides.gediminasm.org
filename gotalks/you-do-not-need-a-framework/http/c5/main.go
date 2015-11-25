package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

func hello(w http.ResponseWriter, req *http.Request) error {
	if errMsg := req.URL.Query().Get("err"); len(errMsg) > 0 {
		return fmt.Errorf(errMsg)
	}
	fmt.Fprintf(w, "hello world\n")
	return nil
}

func main() {
	http.Handle("/hello", Standard.Append(Method("GET")).Then(HandleErr(hello)))
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type ErrorHandler func(w http.ResponseWriter, req *http.Request) error // HLbbb

func HandleErr(h ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil { // HLbbb
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, http.StatusText(http.StatusBadRequest)+" - "+err.Error())
		}
	})
}

var Standard = alice.New(Timing())

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
