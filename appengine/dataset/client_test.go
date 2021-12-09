package dataset_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func TestCreateIfNotExists(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../credentials.json")

	cases := []struct {
		name string
		meta bigquery.TableMetadata
	}{
		{"raw", dataset.CommitsTableMeta},
		{"raw", dataset.PullReqsTableMeta},
		{"raw", dataset.PullReqCommitsTableMeta},
		{"raw", dataset.WorkflowRunsTableMeta},
		{"raw", dataset.WorkflowJobsTableMeta},
	}

	for _, c := range cases {
		ctx := context.Background()
		client, err := dataset.New(ctx)
		if err != nil {
			t.Fatalf("new bigquery client: %v", err)
		}

		if err := client.CreateIfNotExists(ctx, c.name, c.meta); err != nil {
			t.Errorf("create if not exists: %v", err)
		}
	}
}
