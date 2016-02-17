package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/godog"
	"github.com/cucumber/gherkin-go"
)

type apiFeature struct {
	resp   *httptest.ResponseRecorder
	api    *API
	dbMock sqlmock.Sqlmock
}

func (a *apiFeature) reset(interface{}) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	a.resp = httptest.NewRecorder()
	a.api = New(db)
	a.dbMock = mock
}

func (a *apiFeature) thereAreEmployees(users *gherkin.DataTable) error {
	rows := sqlmock.NewRows([]string{"firstname", "lastname"})
	for _, row := range users.Rows {
		rows.AddRow(row.Cells[0].Value, row.Cells[0].Value)
	}
	a.dbMock.ExpectQuery("SELECT (.*) FROM users").WillReturnRows(rows)
	return godog.ErrPending
}

func (a *apiFeature) iSendrequestTo(method, endpoint string) (err error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return
	}

	// handle panic
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	switch endpoint {
	case "/users":
		a.api.users(a.resp, req)
	default:
		err = fmt.Errorf("unknown endpoint: %s", endpoint)
	}
	return
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}
	return nil
}

func (a *apiFeature) theResponseShouldMatchJSON(body *gherkin.DocString) (err error) {
	var expected, actual []byte
	var data interface{}
	if err = json.Unmarshal([]byte(body.Content), &data); err != nil {
		return
	}
	if expected, err = json.Marshal(data); err != nil {
		return
	}
	actual = a.resp.Body.Bytes()
	if !bytes.Equal(actual, expected) {
		err = fmt.Errorf("expected json, does not match actual: %s", string(actual))
	}
	return
}

func apiContext(s *godog.Suite) {
	api := &apiFeature{}

	s.BeforeScenario(api.reset)

	s.Step(`^there are employees:$`, api.thereAreEmployees)
	s.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendrequestTo)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
}
