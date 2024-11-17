package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/auth"
	"net/http"
)

// HandlerOption if code != 0 or err != nil, the controller will abort the request
type HandlerOption func(c *gin.Context) (code int, err error)

type Controller struct {
	actions      map[string]IAction
	name         string
	actionAuthFn ActionAuthFunc
}

func (c *Controller) Name() string {
	return c.name
}

func (c *Controller) AddAction(action IAction) {
	c.actions[action.Name()] = action
}

func (c *Controller) RegisterActions(r gin.IRouter, opts ...HandlerOption) {
	for _, action := range c.actions {
		// register the handler with the router
		r.Any(c.ActionPath(action), c.Handler(action, opts...))
		// register the permissions with the auth middleware
		auth.SetAllowedRoles(c.ActionPath(action), action.AuthorizedRoles()...)
	}
}

func (c *Controller) ActionPath(action IAction) string {
	if v := action.VanityPath(); v != "" {
		return v
	}
	
	return "/" + c.Name() + "/" + action.Name()
}

func (c *Controller) Handler(action IAction, opts ...HandlerOption) gin.HandlerFunc {
	isAuthorized := c.actionAuthFn

	// ensure that an authorization function is always set for the action handler
	if isAuthorized == nil {
		isAuthorized = defaultIsAuthorized(action.AuthorizedRoles()...)
	}

	return func(ctx *gin.Context) {
		if !action.HasVerb(ctx.Request.Method) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if !isAuthorized(ctx) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Set(auth.IsAuthorizedContextKey, true)

		for _, o := range opts {
			if code, err := o(ctx); code != 0 || err != nil {
				if code < 1 {
					code = http.StatusInternalServerError
				}

				ctx.AbortWithStatus(code)
				return
			}
		}

		// call the action handler
		action.GetFn()(ctx)
	}
}

var (
	controllers map[string][]*Controller = make(map[string][]*Controller)
)

func Register(group string, cont *Controller) {
	if _, ok := controllers[group]; !ok {
		controllers[group] = make([]*Controller, 0)
	}

	controllers[group] = append(controllers[group], cont)
}

func GetControllers(group string) []*Controller {
	if cv, ok := controllers[group]; ok {
		return cv
	}

	return nil
}

func GetAllControllers() map[string][]*Controller {
	return controllers
}

func NewController(name string, actions ...map[string]IAction) *Controller {
	if len(actions) < 1 {
		actions = append(actions, make(map[string]IAction))
	}

	return &Controller{
		actions:      actions[0],
		name:         name,
		actionAuthFn: nil,
	}
}
