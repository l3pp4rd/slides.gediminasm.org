package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/icrowley/fake"
)

var build = flag.Bool("build", false, "build database json to stdout")

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Company   string `json:"company"`
	Country   string `json:"country"`
}

var users []*User
var firstname []*User

func main() {
	flag.Parse()

	started := time.Now()
	if *build {
		for i := 0; i < 1000000; i++ {
			users = append(users, &User{
				ID:        i + 1,
				Firstname: fake.FirstName(),
				Lastname:  fake.LastName(),
				Company:   fake.Company(),
				Country:   fake.Country(),
			})
		}
		if err := json.NewEncoder(os.Stdout).Encode(users); err != nil {
			panic(err)
		}
		fmt.Println("built", len(users), "in:", time.Since(started))
		return
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		panic("expecting database from STDIN")
	}

	if err := json.NewDecoder(os.Stdin).Decode(&users); err != nil {
		panic(err)
	}
	fmt.Println("loaded", len(users), "users up into memory in:", time.Since(started))

	go indexing()
	http.Handle("/", timing(http.HandlerFunc(handler)))

	fmt.Println("listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	iparam := func(s string, dflt int) int {
		i, err := strconv.Atoi(s)
		if err != nil || i <= 0 {
			return dflt
		}
		return i
	}

	page := iparam(req.URL.Query().Get("page"), 1) - 1
	limit := iparam(req.URL.Query().Get("limit"), 10)

	lookup := users
	if srch := req.URL.Query().Get("search"); len(srch) > 0 {
		lookup = search(srch)
	}

	from := page * limit
	to := from + limit

	var paged []*User
	switch {
	case len(lookup) >= to:
		paged = lookup[from:to]
	case len(lookup) > from && len(lookup) < to:
		paged = lookup[from:]
	case len(lookup) < from:
		paged = make([]*User, 0)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(paged)
}

func timing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		started := time.Now()
		h.ServeHTTP(w, req)
		fmt.Printf("%s - %s in %s\n", req.Method, req.URL.String(), time.Since(started))
	})
}

func indexing() {
	started := time.Now()

	firstname = make([]*User, len(users))
	copy(firstname, users)
	sort.Sort(byFirstName(firstname))

	fmt.Println("indexed by firstname in:", time.Since(started))
}

func search(pat string) (ret []*User) {
	idx := sort.Search(len(firstname), func(i int) bool {
		return firstname[i].Firstname >= pat
	})
	for len(firstname) > idx && firstname[idx].Firstname == pat {
		ret = append(ret, firstname[idx])
		idx++
	}
	return
}

type byFirstName []*User

func (a byFirstName) Len() int           { return len(a) }
func (a byFirstName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFirstName) Less(i, j int) bool { return a[i].Firstname < a[j].Firstname }
