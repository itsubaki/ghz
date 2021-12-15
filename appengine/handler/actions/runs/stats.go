package runs

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

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.WorkflowRunStatsTableMeta); err != nil {
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

		items := make([]interface{}, 0)
		for i := range stats {
			items = append(items, stats[i])
		}

		if err := dataset.Insert(ctx, datasetName, dataset.WorkflowRunStatsTableMeta.Name, items); err != nil {
			log.Printf("insert items: %v", err)
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
		"select workflow_id, workflow_name, run_id, conclusion, created_at, updated_at, head_sha from `%v` where created_at >= \"%v\" and created_at < \"%v\"",
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
			HeadSHA:      values[6].(string),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func GetStats(owner, repository string, start, end time.Time, list []dataset.WorkflowRun) []dataset.WorkflowRunStats {
	runs := make(map[int64][]dataset.WorkflowRun)
	for _, r := range list {
		v, ok := runs[r.WorkflowID]
		if !ok {
			runs[r.WorkflowID] = make([]dataset.WorkflowRun, 0)
		}

		runs[r.WorkflowID] = append(v, r)
	}

	out := make([]dataset.WorkflowRunStats, 0)
	for _, v := range runs {
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

			sum += math.Pow(r.UpdatedAt.Sub(r.CreatedAt).Minutes()-avg, 2.0)
		}
		variant := sum / count

		var leadtime float64
		{
			list := make([]time.Duration, 0)
			for _, r := range v {
				lt, err := LeadTime(owner, repository, r)
				if err != nil {
					panic(fmt.Sprintf("lead time(%v/%v): %v", owner, repository, err))
				}

				list = append(list, lt...)
			}

			var sum float64
			for _, t := range list {
				sum += t.Minutes()
			}

			if len(list) > 0 {
				leadtime = sum / float64(len(list))
			}
		}

		y, w := start.ISOWeek()
		out = append(out, dataset.WorkflowRunStats{
			Owner:        owner,
			Repository:   repository,
			WorkflowID:   v[0].WorkflowID,
			WorkflowName: v[0].WorkflowName,
			Year:         int64(y),
			Week:         int64(w),
			Start:        civil.DateOf(start),
			End:          civil.DateOf(end),
			RunsPerDay:   runsperday,
			FailureRate:  rate,
			DurationAvg:  avg,
			DurationVar:  variant,
			LeadTime:     leadtime,
		})
	}

	return out
}

func LeadTime(owner, repository string, run dataset.WorkflowRun) ([]time.Duration, error) {
	if run.Conclusion != "success" {
		return make([]time.Duration, 0), nil
	}

	ctx := context.Background()
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}
	datasetName := dataset.Name(owner, repository)

	var id int64
	{
		table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
		query := fmt.Sprintf("select id from `%v` where merge_commit_sha = \"%v\"", table, run.HeadSHA)

		if err := client.Query(ctx, query, func(values []bigquery.Value) {
			if len(values) != 1 {
				return
			}

			if values[0] == nil {
				return
			}

			id = values[0].(int64)
		}); err != nil {
			return nil, fmt.Errorf("query(%v): %v", query, err)
		}

		if id == 0 {
			return make([]time.Duration, 0), nil
		}
	}

	out := make([]time.Duration, 0)
	{
		table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqCommitsTableMeta.Name)
		query := fmt.Sprintf("select date from `%v` where id = %v", table, id)

		if err := client.Query(ctx, query, func(values []bigquery.Value) {
			if len(values) != 1 {
				return
			}

			if values[0] == nil {
				return
			}

			date := values[0].(time.Time)
			out = append(out, run.UpdatedAt.Sub(date))
		}); err != nil {
			return nil, fmt.Errorf("query(%v): %v", query, err)
		}
	}

	return out, nil
}

func NextTime(ctx context.Context, datasetName string) (time.Time, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return time.Now(), fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowRunStatsTableMeta.Name)
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
