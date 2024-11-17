package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/auth"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

// GetApiKeyFunc retrieves the api key from storage. return an error with text 'not found' to correctly handle
// the difference between a storage error and an invalid api key
type GetApiKeyFunc func(ctx context.Context, apiKey string) (auth.ApiKey, error)

func ApiKeyAuthRequired(getApiKey GetApiKeyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(AuthHeader)

		if apiKey == "" {
			if k, ok := c.GetQuery(AuthQueryParam); ok && k != "" {
				apiKey = k
			}
		}

		if apiKey == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "missing authentication credentials"})
			return
		}

		// check if api key is valid
		key, kErr := getApiKey(c, apiKey)

		if kErr != nil {
			if kErr.Error() == "not found" {
				_ = c.AbortWithError(http.StatusForbidden, kErr)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, kErr)
			}

			return
		}

		if !key.GetIsActive() {
			_ = c.AbortWithError(http.StatusForbidden, errors.New("using deactivated api key "+key.GetId()))
			return
		}

		// TODO: verify signature

		user := key.GetUser()

		// add attributes to span
		tSpan := trace.SpanFromContext(c)

		tSpan.SetAttributes(attribute.String("apiKey.id", key.GetId()))

		if user != nil {
			tSpan.SetAttributes(attribute.String("user.id", user.GetId()))
		}

		// add the roles that have been assigned to the ApiKey to the context
		c.Set(auth.RoleContextKey, auth.NewRoleFromApiKeyRole(key.GetRoles()...))

		// add the user to the context for down stream handlers to reference
		c.Set(ContextUserKey, user)

		// add the api key id
		c.Set(ContextApiKeyId, key.GetId())

		c.Next()
	}
}
