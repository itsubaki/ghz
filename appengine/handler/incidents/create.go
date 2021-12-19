package incidents

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func Create(c *gin.Context) {
	var in dataset.Incident
	if err := c.BindJSON(&in); err != nil {
		log.Printf("bind json: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	in.Owner = c.Param("owner")
	in.Repository = c.Param("repository")
	in.ID = NewRandomID()

	if in.PullReqID == 0 {
		log.Printf("pullreq_id(%v) is invalid", in.PullReqID)
		c.Status(http.StatusBadRequest)
		return
	}

	if in.CreatedAt.Year() == 1 {
		log.Printf("created_at(%v) is invalid", in.CreatedAt)
		c.Status(http.StatusBadRequest)
		return
	}

	if in.ResolvedAt.Year() == 1 {
		log.Printf("resolved_at(%v) is invalid", in.ResolvedAt)
		c.Status(http.StatusBadRequest)
		return
	}

	// SELECT count(id) FROM `projectID.datasetName.pullreqs` where id = in.PullRqID
	// if 0 then error

	// datasetName := dataset.Name(in.Owner, in.Repository)

	// ctx := context.Background()
	// if err := dataset.CreateIfNotExists(ctx, datasetName, []bigquery.TableMetadata{
	// 	dataset.IncidentMeta,
	// }); err != nil {
	// 	log.Printf("create if not exists: %v", err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return
	// }

	// items := make([]interface{}, 0)
	// items = append(items, in)

	// if err := dataset.Insert(ctx, datasetName, dataset.IncidentMeta.Name, items); err != nil {
	// 	log.Printf("insert items: %v", err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return
	// }

	c.JSON(http.StatusOK, in)
}
