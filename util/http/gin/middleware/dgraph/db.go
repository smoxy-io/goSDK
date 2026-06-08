package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/db/dgraph"
	"net/http"
)

// Dgraph middleware that configures dgraph database connection management per request
func Dgraph() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := dgraph.BackgroundContext(c)

		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Next()
	}
}
