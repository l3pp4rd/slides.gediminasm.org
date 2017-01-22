package main

import (
	"fmt"
	"log"
	"net/http"
)

func anything(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "handled")
}

func main() {
	http.HandleFunc("/", anything)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
