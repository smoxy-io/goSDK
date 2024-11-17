package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/http/gin/controllers"
	"strings"
)

func sanitizePath(path string) string {
	if path == "" {
		path = "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path
}

func registerHandlers(r gin.IRouter, ctrls []*controllers.Controller) {
	if len(ctrls) == 0 {
		return
	}

	for _, c := range ctrls {
		c.RegisterActions(r)
	}
}
