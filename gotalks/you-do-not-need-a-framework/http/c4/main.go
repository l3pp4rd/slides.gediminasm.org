package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func product(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "category - %s product %s!\n", ps.ByName("cid"), ps.ByName("pid"))
}

func main() {
	router := httprouter.New()
	router.GET("/category/:cid/product/:pid", product)

	http.Handle("/hello", Standard.Append(Method("GET")).ThenFunc(hello)) // HLaaa
	http.Handle("/", Standard.Then(router))                               // HLaaa

	log.Fatal(http.ListenAndServe(":8080", nil))
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
