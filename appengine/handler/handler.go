package handler

import (
	"log"
	"net/http"
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

	g.Use(SetProjectID)
	g.Use(SetTraceID)

	g.Use(gin.Recovery())
	if gin.IsDebugging() {
		g.Use(gin.Logger())
	}

	g.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		log.Printf("%#v", c.Errors.Errors())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			c.Errors.Last().Meta,
		)
	})

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
	if c.GetHeader("X-Appengine-Cron") != "true" {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": "X-Appengine-Cron header is not set to true",
		})
		return
	}

	c.Next()
}

func SetProjectID(c *gin.Context) {
	c.Set("project_id", ProjectID)
	c.Next()
}

func SetTraceID(c *gin.Context) {
	trace := c.GetHeader("X-Cloud-Trace-Context")
	if trace == "" {
		c.Next()
		return
	}

	traceID := strings.Split(trace, "/")[0]
	c.Set("trace_id", traceID)
	c.Next()
}
