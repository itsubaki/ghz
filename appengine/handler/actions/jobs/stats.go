package jobs

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

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.WorkflowJobStatsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	next, err := NextTime(ctx, datasetName)
	if err != nil {
		log.Printf("get lastID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, next)

	for _, d := range calendar.LastNWeeks(52) { // 1 year ~= 52 weeks
		if d.Start.Before(next) {
			// already done it
			continue
		}

		jobs, err := GetJobsWith(ctx, datasetName, d.Start, d.End)
		if err != nil {
			log.Printf("get runs with(%v, %v): %v", d.Start, d.End, err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(jobs) == 0 {
			continue
		}

		stats := GetStats(owner, repository, d.Start, d.End, jobs)

		items := make([]interface{}, 0)
		for i := range stats {
			items = append(items, stats[i])
		}

		if err := dataset.Insert(ctx, datasetName, dataset.WorkflowJobStatsTableMeta.Name, items); err != nil {
			log.Printf("insert items: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
}

func GetStats(owner, repository string, start, end time.Time, list []dataset.WorkflowJob) []dataset.WorkflowJobStats {
	jobs := make(map[string][]dataset.WorkflowJob)
	for _, r := range list {
		v, ok := jobs[r.JobName]
		if !ok {
			jobs[r.JobName] = make([]dataset.WorkflowJob, 0)
		}

		jobs[r.JobName] = append(v, r)
	}

	out := make([]dataset.WorkflowJobStats, 0)
	for k, v := range jobs {
		var failure float64
		var duration time.Duration
		for _, r := range v {
			if r.Conclusion == "failure" {
				failure++
			}

			if r.Conclusion == "success" {
				duration += r.CompletedAt.Sub(r.StartedAt)
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

			sum += math.Pow(r.CompletedAt.Sub(r.StartedAt).Minutes()-avg, 2.0)
		}
		variant := sum / count

		out = append(out, dataset.WorkflowJobStats{
			Owner:        owner,
			Repository:   repository,
			WorkflowID:   v[0].WorkflowID,
			WorkflowName: v[0].WorkflowName,
			JobName:      k,
			Start:        civil.DateOf(start),
			End:          civil.DateOf(end),
			RunsPerDay:   runsperday,
			FailureRate:  rate,
			DurationAvg:  avg,
			DurationVar:  variant,
		})
	}

	return out
}

func GetJobsWith(ctx context.Context, datasetName string, start, end time.Time) ([]dataset.WorkflowJob, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowJobsTableMeta.Name)
	query := fmt.Sprintf(
		"select workflow_id, workflow_name, job_name, conclusion, started_at, completed_at from `%v` where started_at >= \"%v\" and completed_at < \"%v\"",
		table,
		start.Format("2006-01-02 15:04:05 UTC"),
		end.Format("2006-01-02 15:04:05 UTC"),
	)

	out := make([]dataset.WorkflowJob, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		out = append(out, dataset.WorkflowJob{
			WorkflowID:   values[0].(int64),
			WorkflowName: values[1].(string),
			JobName:      values[2].(string),
			Conclusion:   values[3].(string),
			StartedAt:    values[4].(time.Time),
			CompletedAt:  values[5].(time.Time),
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

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowJobStatsTableMeta.Name)
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
