package incidents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	c.Status(http.StatusOK)
}
