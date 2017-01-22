package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func docs(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "v1 documentation")
}

func users(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "users: %s", ps.ByName("id"))
}

func main() {
	router := httprouter.New()
	router.GET("/v1", docs)
	router.GET("/v1/users/{id}", users)

	log.Fatal(http.ListenAndServe(":8080", router))
}
