package runs

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/calendar"
)

func Stats(c *gin.Context) {
	ctx := context.Background()
	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	weeks := c.Query("weeks")
	if weeks == "" {
		weeks = "1"
	}

	w, err := strconv.Atoi(weeks)
	if err != nil {
		log.Printf("atoi %v: %v", weeks, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, d := range calendar.LastNWeeks(w) {
		runs, err := GetRunsWith(ctx, datasetName, d.Start, d.End)
		if err != nil {
			log.Printf("get runs with(%v, %v): %v", d.Start, d.End, err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(runs) == 0 {
			continue
		}

		stats := GetStats(owner, repository, d.Start, d.End, runs)
		for k, v := range stats {
			log.Printf("%v(%v ~ %v): %#v", k, d.Start, d.End, v)
		}

		if err := Insert(ctx, datasetName, stats); err != nil {
			log.Printf("insert: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
}

func GetRunsWith(ctx context.Context, datasetName string, start, end time.Time) ([]dataset.WorkflowRun, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowRunsTableMeta.Name)
	query := fmt.Sprintf(
		"select workflow_id, workflow_name, run_id, conclusion, created_at, updated_at from `%v` where created_at >= \"%v\" and created_at < \"%v\"",
		table,
		start.Format("2006-01-02 15:04:05 UTC"),
		end.Format("2006-01-02 15:04:05 UTC"),
	)

	out := make([]dataset.WorkflowRun, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		out = append(out, dataset.WorkflowRun{
			WorkflowID:   values[0].(int64),
			WorkflowName: values[1].(string),
			RunID:        values[2].(int64),
			Conclusion:   values[3].(string),
			CreatedAt:    values[4].(time.Time),
			UpdatedAt:    values[5].(time.Time),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func GetStats(owner, repository string, start, end time.Time, list []dataset.WorkflowRun) map[int64]dataset.WorkflowRunStats {
	runs := make(map[int64][]dataset.WorkflowRun)
	for _, r := range list {
		v, ok := runs[r.WorkflowID]
		if !ok {
			runs[r.WorkflowID] = make([]dataset.WorkflowRun, 0)
		}

		runs[r.WorkflowID] = append(v, r)
	}

	out := make(map[int64]dataset.WorkflowRunStats)
	for k, v := range runs {
		var failure float64
		var duration time.Duration
		for _, r := range v {
			if r.Conclusion == "failure" {
				failure++
			}

			if r.Conclusion == "success" {
				duration += r.UpdatedAt.Sub(r.CreatedAt)
			}
		}

		count := float64(len(v))
		rate := failure / count
		avg := duration.Minutes() / count
		runsperday := count / (end.Sub(start).Hours() / 24)

		var sum float64
		for _, r := range v {
			if r.Conclusion != "success" {
				continue
			}

			sum += math.Pow((r.UpdatedAt.Sub(r.CreatedAt).Minutes() - avg), 2.0)
		}
		variant := sum / count

		out[k] = dataset.WorkflowRunStats{
			Owner:        owner,
			Repository:   repository,
			WorkflowID:   v[0].WorkflowID,
			WorkflowName: v[0].WorkflowName,
			Start:        start,
			End:          end,
			RunsPerDay:   runsperday,
			FailureRate:  rate,
			DurationAvg:  avg,
			DurationVar:  variant,
		}
	}

	return out
}

func Insert(ctx context.Context, datasetName string, stats map[int64]dataset.WorkflowRunStats) error {
	return nil
}
