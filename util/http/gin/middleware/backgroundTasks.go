package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

// BackgroundTasks middleware that provides support for clean shutdowns when background tasks are needed
func BackgroundTasks(wg *sync.WaitGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		// add the server's wait group to the context so it's available downstream
		c.Set(ContextBackgroundTasksWg, wg)

		c.Next()
	}
}

func GetBackgroundTasksWg(ctx context.Context) *sync.WaitGroup {
	wg, ok := ctx.Value(ContextBackgroundTasksWg).(*sync.WaitGroup)

	if !ok {
		return nil
	}

	return wg
}
