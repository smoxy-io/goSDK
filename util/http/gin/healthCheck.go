package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Status string `json:"status"`
		}{Status: "ok"})
	}
}
