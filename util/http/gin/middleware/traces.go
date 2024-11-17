package middleware

import (
	"net/http"
	"slices"
)

var (
	defaultNoTraceEndpoints = []string{
		"/health",
		"/metrics",
	}
)

func FilterTraces(noTraceEndpoints ...string) func(*http.Request) bool {
	if len(noTraceEndpoints) == 0 {
		noTraceEndpoints = defaultNoTraceEndpoints
	}

	return func(req *http.Request) bool {
		return slices.Index(noTraceEndpoints, req.URL.Path) == -1
	}
}
