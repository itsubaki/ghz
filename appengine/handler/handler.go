package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/handler/actions/jobs"
	"github.com/itsubaki/ghz/appengine/handler/actions/runs"
	"github.com/itsubaki/ghz/appengine/handler/commits"
	"github.com/itsubaki/ghz/appengine/handler/events"
	"github.com/itsubaki/ghz/appengine/handler/incidents"
	"github.com/itsubaki/ghz/appengine/handler/pullreqs"
	prcommits "github.com/itsubaki/ghz/appengine/handler/pullreqs/commits"
	"github.com/itsubaki/ghz/appengine/handler/releases"
)

func New() *gin.Engine {
	g := gin.New()

	g.Use(SetTraceID)

	g.Use(gin.Recovery())
	if gin.IsDebugging() {
		g.Use(gin.Logger())
	}

	Root(g)
	Status(g)
	Fetch(g)
	Incidents(g)

	return g
}

func Root(g *gin.Engine) {
	g.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func Status(g *gin.Engine) {
	g.GET("/status", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func Fetch(g *gin.Engine) {
	r := g.Group("/_fetch")
	r.Use(XAppEngineCron)

	r.GET("/:owner/:repository/_init", Init)
	r.GET("/:owner/:repository/commits", commits.Fetch)
	r.GET("/:owner/:repository/events", events.Fetch)
	r.GET("/:owner/:repository/releases", releases.Fetch)
	r.GET("/:owner/:repository/pullreqs", pullreqs.Fetch)
	r.GET("/:owner/:repository/pullreqs/update", pullreqs.Update)
	r.GET("/:owner/:repository/pullreqs/commits", prcommits.Fetch)
	r.GET("/:owner/:repository/actions/runs", runs.Fetch)
	r.GET("/:owner/:repository/actions/runs/update", runs.Update)
	r.GET("/:owner/:repository/actions/jobs", jobs.Fetch)
	r.GET("/:owner/:repository/actions/jobs/update", jobs.Update)
}

func Incidents(g *gin.Engine) {
	r := g.Group("/incidents")

	r.POST("/:owner/:repository", incidents.Create)
}

func XAppEngineCron(c *gin.Context) {
	// https://cloud.google.com/appengine/docs/standard/go/scheduling-jobs-with-cron-yaml
	// Requests from the Cron Service will contain the following HTTP header:
	// X-Appengine-Cron: true
	if c.GetHeader("X-Appengine-Cron") != "true" {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": "X-Appengine-Cron header is not set to true",
		})
		return
	}

	c.Next()
}

func SetTraceID(c *gin.Context) {
	value := c.GetHeader("X-Cloud-Trace-Context")
	if value == "" {
		c.Next()
		return
	}

	// https://cloud.google.com/trace/docs/setup
	// The header specification is:
	// "X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=TRACE_TRUE"
	ids := strings.Split(strings.Split(value, ";")[0], "/")
	c.Set("trace_id", ids[0])

	// https://cloud.google.com/trace/docs/setup
	// SPAN_ID is the decimal representation of the (unsigned) span ID.
	i, err := strconv.ParseUint(ids[1], 10, 64)
	if err != nil {
		log.Printf("parse %v: %v", ids[1], err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md#retrieving-the-traceid-and-spanid
	// MUST be a 16-hex-character lowercase string
	c.Set("span_id", fmt.Sprintf("%016x", i))

	c.Set("trace_true", false)
	if len(strings.Split(value, ";")) > 1 && strings.Split(value, ";")[1] == "o=1" {
		c.Set("trace_true", true)
	}

	c.Next()
}
