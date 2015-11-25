package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

func hello(w http.ResponseWriter, req *http.Request) error {
	switch req.URL.Query().Get("err") {
	case "400":
		return NewHTTPError(400, "some 400 error")
	case "500":
		return NewHTTPError(500, "unexpected error")
	}
	fmt.Fprintf(w, "hello world\n")
	return nil
}

func main() {
	http.Handle("/hello", Standard.Append(Method("GET")).Then(HandleErr(hello)))
	http.Handle("/", http.NotFoundHandler())

	log.Println("listening on :8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (e *HTTPError) JSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	body := struct {
		Status    int    `json:"status"`
		Error     string `json:"error"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    e.statusCode,
		Error:     http.StatusText(e.statusCode),
		Message:   e.message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	response, err := json.Marshal(body)
	if err != nil {
		return err
	}
	w.WriteHeader(e.statusCode)
	_, err = w.Write(response)
	return err
}

type ErrorHandler func(w http.ResponseWriter, req *http.Request) error

func HandleErr(h ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			switch t := err.(type) {
			case *HTTPError:
				t.JSON(w)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, t.Error())
			}
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

type HTTPError struct {
	statusCode int
	message    string
}

func NewHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{code, msg}
}

func (e *HTTPError) Error() string {
	return e.message
}
