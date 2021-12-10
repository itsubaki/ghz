package pullreqs

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func Update(c *gin.Context) {

}

func GetPullReqsWith(ctx context.Context, datasetName, state string) error {
	client, err := dataset.New(ctx)
	if err != nil {
		return fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where state = %v", table, state)

	log.Printf("query: %v", query)
	return nil
}
