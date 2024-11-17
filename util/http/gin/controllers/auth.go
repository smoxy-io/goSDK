package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/auth"
	"slices"
)

func defaultIsAuthorized(allowedRoles ...auth.Role) ActionAuthFunc {
	if len(allowedRoles) < 1 {
		// no roles are allowed
		// deny all by default
		return func(c *gin.Context) bool {
			return false
		}
	}

	if slices.Index(allowedRoles, auth.RoleAnonymous) != -1 {
		// anonymous role is allowed, so all roles are allowed
		return func(c *gin.Context) bool {
			return true
		}
	}

	return func(c *gin.Context) bool {
		roles := auth.GetRolesFromCtx(c)

		if roles == auth.RoleAnonymous {
			// anonymous role is not allowed here
			return false
		}

		for _, role := range allowedRoles {
			if roles.Has(role) {
				return true
			}
		}

		return false
	}
}
