package pullreqs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/beta"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func GetStats(c *gin.Context) {
	ctx := beta.Context()

	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)
	cachekey := fmt.Sprintf("stats_pullreq_%v", datasetName)

	cache, err := beta.MemGet(ctx, cachekey)
	if err != nil && err != beta.ErrCacheMiss {
		log.Printf("memcache get: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if err == nil {
		c.String(http.StatusOK, string(cache))
		return
	}

	client, err := dataset.New(ctx)
	if err != nil {
		log.Printf("new bigquery client: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqStatsTableMeta.Name)
	query := fmt.Sprintf("select year, week, start, `end`, merged_per_day, duration_avg, duration_var from `%v`", table)
	out := make([]dataset.PullReqStats, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		out = append(out, dataset.PullReqStats{
			Owner:        owner,
			Repository:   repository,
			Year:         values[0].(int64),
			Week:         values[1].(int64),
			Start:        values[2].(civil.Date),
			End:          values[3].(civil.Date),
			MergedPerDay: values[4].(float64),
			DurationAvg:  values[5].(float64),
			DurationVar:  values[6].(float64),
		})
	}); err != nil {
		log.Printf("query(%v): %v", query, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	if err := beta.MemSet(ctx, cachekey, b, 24*time.Hour); err != nil {
		log.Printf("memcaceh set(%v): %v", cachekey, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, out)
}
