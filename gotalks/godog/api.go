package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type user struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{db: db}
}

func (a *API) readUsers() ([]*user, error) {
	users := make([]*user, 0) // non nil slice
	rows, err := a.db.Query("SELECT firstname, lastname FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &user{}
		if err := rows.Scan(&user.Firstname, &user.Lastname); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (a *API) users(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	users, err := a.readUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read users from db: %s", err), 500)
		return
	}

	response := struct {
		Users []*user `json:"users"`
	}{Users: users}

	data, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal users: %s", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (a *API) auth(w http.ResponseWriter, req *http.Request) {
	if _, hasHeader := req.Header["Authorization"]; !hasHeader {
		http.Error(w, "authentication required", 401)
		return
	}
	auth := strings.SplitN(req.Header["Authorization"][0], " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "bad syntax", http.StatusBadRequest)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 || pair[0] != "user" || pair[1] != "pass" {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	a.users(w, req)
}
