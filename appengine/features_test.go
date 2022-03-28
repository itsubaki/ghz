package main_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/handler"
	"github.com/jfilipczyk/gomatch"
)

var api = &apiFeature{}

type apiFeature struct {
	header http.Header
	body   io.Reader
	resp   *httptest.ResponseRecorder

	server *gin.Engine
	keep   map[string]interface{}
	result [][]bigquery.Value
}

func (a *apiFeature) start() {
	a.server = handler.New()
	a.keep = make(map[string]interface{})
	a.result = make([][]bigquery.Value, 0)
}

func (a *apiFeature) reset(sc *godog.Scenario) {
	a.header = make(http.Header)
	a.body = nil
	a.resp = httptest.NewRecorder()
	a.result = make([][]bigquery.Value, 0)
}

func (a *apiFeature) replace(str string) string {
	for k, v := range a.keep {
		switch val := v.(type) {
		case string:
			str = strings.Replace(str, k, val, -1)
		}
	}

	return str
}

func (a *apiFeature) Request(method, endpoint string) error {
	r := a.replace(endpoint)
	req := httptest.NewRequest(method, r, a.body)
	req.Header = a.header

	a.server.ServeHTTP(a.resp, req)
	return nil
}

func (a *apiFeature) ResponseCodeShouldBe(code int) error {
	if code == a.resp.Code {
		return nil
	}

	return fmt.Errorf("got=%v, want=%v", a.resp.Code, code)
}

func (a *apiFeature) ResponseShouldMatchJSON(body *godog.DocString) error {
	want := a.replace(body.Content)
	got := a.resp.Body.String()

	ok, err := gomatch.NewDefaultJSONMatcher().Match(want, got)
	if err != nil {
		return fmt.Errorf("got=%v, want=%v, match: %v", got, want, err)
	}

	if !ok {
		return fmt.Errorf("got=%v, want=%v", got, want)
	}

	return nil
}

func (a *apiFeature) SetHeader(k, v string) error {
	a.header.Add(k, v)
	return nil
}

func (a *apiFeature) SetRequestBody(body *godog.DocString) error {
	r := a.replace(body.Content)
	a.body = bytes.NewBuffer([]byte(r))
	return nil
}

func (a *apiFeature) IncidentsExists(incidents *godog.Table) error {
	owner := incidents.Rows[1].Cells[0].Value
	repository := incidents.Rows[1].Cells[1].Value
	dsn := dataset.Name(owner, repository)

	items := make([]interface{}, 0)
	for i := 1; i < len(incidents.Rows); i++ {
		var count int64
		if err := dataset.Query(context.Background(),
			fmt.Sprintf(
				"SELECT count(*) FROM `%v.%v.%v` WHERE sha = \"%v\"",
				dataset.ProjectID, dsn, dataset.IncidentsMeta.Name,
				incidents.Rows[i].Cells[3].Value,
			),
			func(values []bigquery.Value) {
				count = values[0].(int64)
			},
		); err != nil {
			return fmt.Errorf("query: %v", err)
		}

		if count > 0 {
			// already exists
			continue
		}

		resolvedAt, err := time.Parse("2006-01-02 15:04:05 UTC", incidents.Rows[i].Cells[4].Value)
		if err != nil {
			return fmt.Errorf("parse(%v): %v", incidents.Rows[i].Cells[4].Value, err)
		}

		items = append(items, dataset.Incident{
			Owner:       incidents.Rows[i].Cells[0].Value,
			Repository:  incidents.Rows[i].Cells[1].Value,
			Description: incidents.Rows[i].Cells[2].Value,
			SHA:         incidents.Rows[i].Cells[3].Value,
			ResolvedAt:  resolvedAt,
		})
	}

	if err := dataset.Insert(context.Background(), dsn, dataset.IncidentsMeta.Name, items); err != nil {
		return fmt.Errorf("insert into %v: %v", dsn, err)
	}

	return nil
}

func (a *apiFeature) ExecuteQuery(query string) error {
	q := strings.ReplaceAll(query, "$PROJECT_ID", dataset.ProjectID)
	if err := dataset.Query(context.Background(), q, func(values []bigquery.Value) {
		a.result = append(a.result, values)
	}); err != nil {
		return fmt.Errorf("execute query: %v", err)
	}

	return nil
}

func (a *apiFeature) QueryResult(result *godog.Table) error {
	for i := range a.result {
		for j := range a.result[i] {
			name := result.Rows[0].Cells[j].Value
			want := result.Rows[i+1].Cells[j].Value
			switch got := a.result[i][j].(type) {
			case string:
				if want != got {
					return fmt.Errorf("%v got=%v, want=%v", name, got, want)
				}
			case civil.Date:
				if want != got.String() {
					return fmt.Errorf("%v got=%v, want=%v", name, got, want)
				}
			case int64:
				p, err := strconv.ParseInt(want, 10, 64)
				if err != nil {
					return fmt.Errorf("parse int(%v): %v", want, err)
				}

				if p != got {
					return fmt.Errorf("%v got=%v, want=%v", name, got, want)
				}
			case float64:
				p, err := strconv.ParseFloat(want, 64)
				if err != nil {
					return fmt.Errorf("parse float(%v): %v", want, err)
				}

				if p != got {
					return fmt.Errorf("%v got=%v, want=%v", name, got, want)
				}
			}
		}
	}

	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		gin.SetMode(gin.ReleaseMode)
		api.start()
	})

	ctx.AfterSuite(func() {})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(api.reset)

	ctx.Step(`^I set "([^"]*)" header with "([^"]*)"$`, api.SetHeader)
	ctx.Step(`^I set request body:$`, api.SetRequestBody)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.Request)
	ctx.Step(`^the response code should be (\d+)$`, api.ResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.ResponseShouldMatchJSON)
	ctx.Step(`^the following incidents exist:$`, api.IncidentsExists)
	ctx.Step(`^I execute query with:$`, api.ExecuteQuery)
	ctx.Step(`^I get the following result:$`, api.QueryResult)
}
