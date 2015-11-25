package main

import (
	"log"
	"net/http"
	"time"
)

func Timing() Chainable {
	return func(cli Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			start := time.Now()
			resp, err := cli.Do(r)
			log.Printf("%s - %s took: %s\n", r.Method, r.URL.String(), time.Since(start))
			return resp, err
		})
	}
}

func main() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/hello", nil)
	New(Timing()).Then(http.DefaultClient).Do(req)
}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type ClientFunc func(*http.Request) (*http.Response, error)

func (f ClientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

type Chainable func(Client) Client

type Chain []Chainable

func New(ch ...Chainable) Chain {
	return Chain(ch)
}

func (c Chain) Then(cli Client) Client {
	final := cli
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}
