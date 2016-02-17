package api

import (
	"database/sql"
	"net/http"
)

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{db: db}
}

func (a *API) users(w http.ResponseWriter, req *http.Request) {

}
