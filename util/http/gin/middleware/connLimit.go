package middleware

import "github.com/gin-gonic/gin"

func MaxConns(n int) gin.HandlerFunc {
	conns := make(chan struct{}, n)

	return func(c *gin.Context) {
		acquire(conns)       // before request
		defer release(conns) // after request
		c.Next()
	}
}

func acquire(c chan<- struct{}) {
	c <- struct{}{}
}

func release(c <-chan struct{}) {
	<-c
}
