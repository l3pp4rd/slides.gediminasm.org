package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func main() {
	http.HandleFunc("/hello", hello)
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
