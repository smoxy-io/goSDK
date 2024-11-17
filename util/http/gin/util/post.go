package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Post interface {
	IsValid() bool
	Sanitize()
}

func UnmarshalAndValidatePost[T Post](c *gin.Context, post T) bool {
	pErr := c.ShouldBindJSON(post)

	if pErr != nil {
		if pErr != io.EOF {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("param binding error: %s", pErr))
			return false
		}
	}

	post.Sanitize()

	if !post.IsValid() {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request %v", post))
		return false
	}

	return true
}
