package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

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
}

func (a *apiFeature) start() {
	a.server = handler.New()
	a.keep = make(map[string]interface{})
}

func (a *apiFeature) reset(sc *godog.Scenario) {
	a.header = make(http.Header)
	a.body = nil
	a.resp = httptest.NewRecorder()
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

func (a *apiFeature) IncidentsExists(incidents *godog.Table) error {
	items := make([]interface{}, 0)
	for i := 1; i < len(incidents.Rows); i++ {
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

	owner := incidents.Rows[1].Cells[0].Value
	repository := incidents.Rows[1].Cells[1].Value
	_, dsn := dataset.Name(owner, repository)

	if err := dataset.Insert(context.Background(), dsn, dataset.IncidentsMeta.Name, items); err != nil {
		return fmt.Errorf("insert into %v: %v", dsn, err)
	}

	return nil
}

func (a *apiFeature) ExecuteQuery(query string) error {
	return godog.ErrPending
}

func (a *apiFeature) QueryResult(result *godog.Table) error {
	return godog.ErrPending
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	id, dsn := dataset.Name("itsubaki", "q")
	dataset.Delete(context.Background(), id, dsn, []string{
		dataset.CommitsMeta.Name,
		dataset.EventsMeta.Name,
		dataset.IncidentsMeta.Name,
		dataset.PullReqCommitsMeta.Name,
		dataset.PullReqsMeta.Name,
		dataset.ReleasesMeta.Name,
		dataset.WorkflowRunsMeta.Name,
		dataset.WorkflowJobsMeta.Name,
	})

	ctx.BeforeSuite(func() {
		gin.SetMode(gin.ReleaseMode)
		api.start()
	})

	ctx.AfterSuite(func() {})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(api.reset)

	ctx.Step(`^I set "([^"]*)" header with "([^"]*)"$`, api.SetHeader)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.Request)
	ctx.Step(`^the response code should be (\d+)$`, api.ResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.ResponseShouldMatchJSON)
	ctx.Step(`^the following incidents exist:$`, api.IncidentsExists)
	ctx.Step(`^I execute query with:$`, api.ExecuteQuery)
	ctx.Step(`^I get the following result:$`, api.QueryResult)
}
