package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	RecoveryHtmlTmplName = "error.tmpl"
)

type RecoveryData struct {
	Error any `json:"error,omitempty"`
	Data  any `json:"data,omitempty"`
}

type RecoveryHandlerFunc = func(c *gin.Context, err any)

func Recovery(f RecoveryHandlerFunc) gin.HandlerFunc {
	return RecoveryWithWriter(f, gin.DefaultErrorWriter)
}

func RecoveryWithWriter(f RecoveryHandlerFunc, out io.Writer) gin.HandlerFunc {
	var logger *log.Logger

	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}

	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				if logger != nil {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					goErr := errors.Wrap(e, 3)
					reset := string([]byte{27, 91, 48, 109})
					logger.Printf("[Recovery] panic recovered:\n\n%s%s\n\n%s%s", httpRequest, goErr.Error(), goErr.Stack(), reset)
				}

				f(c, e)
			}
		}()
		c.Next() // execute all the handlers
	}
}

func RecoveryHandlerJson(content any, code int) RecoveryHandlerFunc {
	if code < 500 || code > 599 {
		code = http.StatusInternalServerError
	}

	return func(c *gin.Context, err any) {
		data := RecoveryData{
			Data:  content,
			Error: err,
		}

		c.JSON(code, data)
	}
}

func RecoveryHandlerHtml(content any, code int) RecoveryHandlerFunc {
	if code < 500 || code > 599 {
		code = http.StatusInternalServerError
	}

	return func(c *gin.Context, err any) {
		data := RecoveryData{
			Data:  content,
			Error: err,
		}

		c.HTML(code, RecoveryHtmlTmplName, data)
	}
}

// DefaultRecoveryHandler creates a default "json" or "html" recovery handler
//
// contentType must be "json" or "html" (panics on invalid value)
// content will not change between panic recoveries
// statusCode must be a 5xx http response code (default: http.StatusInternalServerError)
//
// the template for the "html" recovery handler MUST be loaded in gin (its path
// must match the value of RecoveryHtmlTmplName) and a RecoveryData object will
// be passed to the template
// ```go
// r := gin.New()
// r.Use(Recovery(DefaultRecoveryHandler("html", gin.H{"foo":"bar"}, 500)))
// r.LoadHTMLFiles(RecoveryHtmlTmplName)
// ```
func DefaultRecoveryHandler(contentType string, content any, statusCode ...int) RecoveryHandlerFunc {
	code := http.StatusInternalServerError

	if len(statusCode) > 0 && statusCode[0] > 499 && statusCode[0] < 600 {
		code = statusCode[0]
	}

	switch strings.ToLower(contentType) {
	case "json", "application/json":
		return RecoveryHandlerJson(content, code)
	case "html", "text/html":
		return RecoveryHandlerHtml(content, code)
	default:
		panic("unsupported recovery handler response type. must be one of: json, html")
	}
}
