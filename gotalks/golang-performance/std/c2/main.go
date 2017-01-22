package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Mux struct {
	// references to services
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const users = "/v1/users/"

	switch {
	case req.URL.Path == "/v1":
		m.docs(w, req)
	case strings.Index(req.URL.Path, users) == 0:
		m.users(w, req, req.URL.Path[len(users):])
	default:
		http.NotFoundHandler().ServeHTTP(w, req)
	}
}

func (m *Mux) docs(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "v1 documentation")
}

func (m *Mux) users(w http.ResponseWriter, req *http.Request, id string) {
	fmt.Fprintf(w, "users: %s", id)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", &Mux{}))
}
