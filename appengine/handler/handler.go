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
	f := g.Group("/_fetch")
	f.Use(XAppEngineCron)

	f.GET("/:owner/:repository/commits", commits.Fetch)
	f.GET("/:owner/:repository/pullreqs", pullreqs.Fetch)
	f.GET("/:owner/:repository/pullreqs/update", pullreqs.Update)
	f.GET("/:owner/:repository/pullreqs/commits", prcommits.Fetch)
	f.GET("/:owner/:repository/actions/runs", runs.Fetch)
	f.GET("/:owner/:repository/actions/jobs", jobs.Fetch)
}
