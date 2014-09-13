package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func job(name string) {
	time.Sleep(time.Second * 10)
	log.Printf("Job %s was handled successfully\n", name)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/job/{name}", func(w http.ResponseWriter, req *http.Request) {
		go job(mux.Vars(req)["name"])
		fmt.Fprintf(w, "Job is now being processed!")
	})
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
