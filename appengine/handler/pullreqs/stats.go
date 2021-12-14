package pullreqs

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/calendar"
)

func Stats(c *gin.Context) {
	ctx := context.Background()
	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.PullReqStatsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	next, err := NextTime(ctx, datasetName)
	if err != nil {
		log.Printf("next: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, next)

	for _, d := range calendar.LastNWeeks(52) { // 1 year ~= 52 weeks
		if d.Start.Before(next) {
			// already done it
			continue
		}

		prs, err := GetPullReqsWith(ctx, datasetName, d.Start, d.End)
		if err != nil {
			log.Printf("get pullreqs: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(prs) == 0 {
			continue
		}

		var sumavg, sumvar float64
		for _, pr := range prs {
			commits, err := GetPullReqCommits(ctx, datasetName, pr.ID)
			if err != nil {
				log.Printf("get pullreq commits: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}

			if len(commits) == 0 {
				continue
			}

			var sum float64
			for _, c := range commits {
				sum += pr.MergedAt.Sub(c.Date).Minutes()
			}

			count := float64(len(commits))
			avg := sum / count
			sumavg += avg

			var sum2 float64
			for _, c := range commits {
				sum2 += math.Pow(pr.MergedAt.Sub(c.Date).Minutes()-avg, 2.0)
			}

			variant := sum2 / count
			sumvar += variant
		}

		count := float64(len(prs))
		perday := count / (d.End.Sub(d.Start).Hours() / 24)
		duravg := sumavg / count
		durvar := sumvar / count
		y, w := d.Start.ISOWeek()

		items := make([]interface{}, 0)
		items = append(items, dataset.PullReqStats{
			Owner:        owner,
			Repository:   repository,
			Year:         int64(y),
			Week:         int64(w),
			Start:        civil.DateOf(d.Start),
			End:          civil.DateOf(d.End),
			MergedPerDay: perday,
			DurationAvg:  duravg,
			DurationVar:  durvar,
		})

		if err := dataset.Insert(ctx, datasetName, dataset.PullReqStatsTableMeta.Name, items); err != nil {
			log.Printf("insert items(%v): %v", items, err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}

func GetPullReqCommits(ctx context.Context, datasetName string, id int64) ([]dataset.PullReqCommits, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqCommitsTableMeta.Name)
	query := fmt.Sprintf(
		"select date from `%v` where id = %v",
		table,
		id,
	)

	out := make([]dataset.PullReqCommits, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		out = append(out, dataset.PullReqCommits{
			Date: values[0].(time.Time),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func GetPullReqsWith(ctx context.Context, datasetName string, start, end time.Time) ([]dataset.PullReqs, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
	query := fmt.Sprintf(
		"select id, merged_at from `%v` where merged_at >= \"%v\" and merged_at < \"%v\"",
		table,
		start.Format("2006-01-02 15:04:05 UTC"),
		end.Format("2006-01-02 15:04:05 UTC"),
	)

	out := make([]dataset.PullReqs, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		out = append(out, dataset.PullReqs{
			ID:       values[0].(int64),
			MergedAt: values[1].(time.Time),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func NextTime(ctx context.Context, datasetName string) (time.Time, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return time.Now(), fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqStatsTableMeta.Name)
	query := fmt.Sprintf("select max(start) from `%v` limit 1", table)

	var out time.Time
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		date := values[0].(civil.Date)
		out = time.Date(date.Year, date.Month, date.Day+1, 0, 0, 0, 0, time.UTC)
	}); err != nil {
		return time.Now(), fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}
