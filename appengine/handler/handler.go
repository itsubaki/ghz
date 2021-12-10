package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/handler/actions/jobs"
	"github.com/itsubaki/ghstats/appengine/handler/actions/runs"
	"github.com/itsubaki/ghstats/appengine/handler/commits"
	"github.com/itsubaki/ghstats/appengine/handler/pullreqs"
	prcommits "github.com/itsubaki/ghstats/appengine/handler/pullreqs/commits"
)

func New() *gin.Engine {
	g := gin.New()

	g.Use(gin.Recovery())
	if gin.IsDebugging() {
		g.Use(gin.Logger())
	}

	Root(g)
	Status(g)
	Fetch(g)
	Stats(g)

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

func XAppEngineCron(c *gin.Context) {
	if c.GetHeader("X-Appengine-Cron") != "true" {
		log.Printf("X-Appengine-Cron header is not set to true")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Next()
}

func Fetch(g *gin.Engine) {
	f := g.Group("/fetch")
	f.Use(XAppEngineCron)

	f.GET("/commits", commits.Fetch)
	f.GET("/pullreqs", pullreqs.Fetch)
	f.GET("/pullreqs/update", pullreqs.Update)
	f.GET("/pullreqs/commits", prcommits.Fetch)
	f.GET("/actions/runs", runs.Fetch)
	f.GET("/actions/jobs", jobs.Fetch)
}

func Stats(g *gin.Engine) {
	s := g.Group("/stats")
	s.Use(XAppEngineCron)

	s.GET("/pullreqs", pullreqs.Stats)
	s.GET("/actions/runs", runs.Stats)
	s.GET("/actions/jobs", jobs.Stats)
}
