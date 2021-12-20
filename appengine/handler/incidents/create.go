package incidents

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/speps/go-hashids"
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

	if in.CreatedAt.Year() == 1 {
		message := fmt.Sprintf("created_at(%v) is invalid", in.CreatedAt)
		log.Printf(message)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
		})
		return
	}

	if in.ResolvedAt.Year() == 1 {
		message := fmt.Sprintf("resolved_at(%v) is invalid", in.ResolvedAt)
		log.Printf(message)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
		})
		return
	}

	ctx := context.Background()
	datasetName := dataset.Name(in.Owner, in.Repository)
	exists, err := Exists(ctx, datasetName, in.PullReqNumber)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !exists {
		message := fmt.Sprintf("pullreq number(%v) is not exists", in.PullReqNumber)
		log.Printf(message)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
		})
		return
	}

	if err := dataset.CreateIfNotExists(ctx, datasetName, []bigquery.TableMetadata{
		dataset.IncidentMeta,
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	items := make([]interface{}, 0)
	items = append(items, in)

	if err := dataset.Insert(ctx, datasetName, dataset.IncidentMeta.Name, items); err != nil {
		log.Printf("insert items: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, in)
}

func Exists(ctx context.Context, datasetName string, number int64) (bool, error) {
	client := dataset.New(ctx)
	defer client.Close()

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("select count(number) from `%v` where number = %v", table, number)

	var count int64
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		count = values[0].(int64)
	}); err != nil {
		return false, fmt.Errorf("query(%v): %v", query, err)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func NewRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return NewID(11, strconv.Itoa(rand.Int()))
}

func NewID(digit int, seed ...string) string {
	if digit == 1 {
		panic(fmt.Sprintf("digit=%d. digit must be greater than 1", digit))
	}

	hd := hashids.NewData()
	hd.MinLength = digit
	hd.Salt = strings.Join(seed, "")

	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	id, err := h.Encode([]int{42})
	if err != nil {
		panic(err)
	}

	return id
}
