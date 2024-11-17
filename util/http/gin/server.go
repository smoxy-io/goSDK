package gin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/http/gin/controllers"
	"github.com/smoxy-io/goSDK/util/http/gin/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	srv              *gin.Engine
	httpSrv          *http.Server
	telemId          string
	backgroundWg     *sync.WaitGroup
	healthCheckRoute string
	connLimit        int
	middleware       map[string][]gin.HandlerFunc
	recoveryHandler  middleware.RecoveryHandlerFunc
	noTraceEndpoints []string
}

var (
	server *Server
)

// NewServer create a new HTTP server that uses the controller/action pattern
func NewServer(telemId string) *Server {
	if server != nil {
		return server
	}

	server = &Server{
		srv:              nil,
		httpSrv:          nil,
		telemId:          strings.TrimSpace(telemId),
		backgroundWg:     nil,
		healthCheckRoute: "/health",
		connLimit:        0,
		middleware:       make(map[string][]gin.HandlerFunc),
		recoveryHandler:  middleware.DefaultRecoveryHandler("json", nil, 500),
		noTraceEndpoints: make([]string, 0),
	}

	server.middleware["main"] = make([]gin.HandlerFunc, 0)

	return server
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.middleware["main"] = append(s.middleware["main"], middleware...)
}

func (s *Server) Group(path string, middleware ...gin.HandlerFunc) {
	if _, ok := s.middleware[path]; !ok {
		s.middleware[path] = make([]gin.HandlerFunc, 0)
	}

	s.middleware[path] = append(s.middleware[path], middleware...)
}

func (s *Server) WithHealthCheckRoute(route string) *Server {
	s.healthCheckRoute = route
	return s
}

func (s *Server) WithConnLimit(limit int) *Server {
	s.connLimit = limit
	return s
}

func (s *Server) WithRecoveryHandler(recoveryHandler middleware.RecoveryHandlerFunc) *Server {
	s.recoveryHandler = recoveryHandler
	return s
}

func (s *Server) WithNoTraceEndpoints(endpoints ...string) *Server {
	s.noTraceEndpoints = append(s.noTraceEndpoints, endpoints...)
	return s
}

// ListenAndServe non-blocking ListenAndServe function
func (s *Server) ListenAndServe(address string) error {
	if s.srv != nil {
		return nil
	}

	s.backgroundWg = &sync.WaitGroup{}

	s.srv = gin.New()

	// register middleware that is always used
	s.srv.Use(gin.Logger())

	if s.connLimit > 0 {
		// MaxConns middleware needs to be BEFORE the recovery handler so that crashes do not mess up the conn count
		s.srv.Use(middleware.MaxConns(s.connLimit))
	}

	s.srv.Use(middleware.Recovery(s.recoveryHandler))
	s.srv.Use(middleware.BackgroundTasks(s.backgroundWg))

	if s.telemId != "" {
		s.srv.Use(otelgin.Middleware(s.telemId, otelgin.WithFilter(middleware.FilterTraces(s.noTraceEndpoints...))))
	}

	// register middleware
	for _, m := range s.middleware["main"] {
		s.srv.Use(m)
	}

	// register the health check route
	s.srv.GET(s.healthCheckRoute, HealthCheck())
	// TODO: add /metrics handler

	// register controllers
	allCtrls := controllers.GetAllControllers()

	for path, ctrls := range allCtrls {
		path = sanitizePath(path)

		group := s.srv.Group(path)

		if mdlware, ok := s.middleware[path]; ok {
			// register middleware for the group
			for _, m := range mdlware {
				group.Use(m)
			}
		}

		// register handlers for the group
		registerHandlers(group, ctrls)
	}

	// setup the http server
	s.httpSrv = &http.Server{
		Addr:    address,
		Handler: s.srv,
	}

	go func() {
		// blocks until server is stopped
		if err := s.httpSrv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Printf("listening on %s\n", address)

	return nil
}

func (s *Server) Stop() error {
	if s.srv == nil {
		return nil
	}

	done := make(chan error, 1)
	tmr := time.NewTimer(20 * time.Second)

	go func() {
		defer func() {
			close(done)
		}()

		sErr := s.httpSrv.Shutdown(context.Background())

		// wait for any background tasks to finish (like short link metrics)
		s.backgroundWg.Wait()

		done <- sErr
	}()

	var err error

	select {
	case err = <-done:
		// graceful stop returned (succeeded if err == nil)
		// intentionally blank
	case <-tmr.C:
		// graceful stop timed out
		return fmt.Errorf("rest server shutdown timed out")
	}

	s.srv = nil
	s.httpSrv = nil

	return err
}
