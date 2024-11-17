package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/smoxy-io/goSDK/util/auth"
	"net/http"
	"slices"
)

type ActionFunc func(c *gin.Context)
type ActionAuthFunc func(c *gin.Context) bool

type IAction interface {
	HasVerb(verb string) bool
	GetFn() ActionFunc
	AuthorizedRoles() []auth.Role
	Name() string
	WithAuthRoles(roles ...auth.Role) IAction
	WithVanityPath(path string) IAction
	VanityPath() string
}

type Action struct {
	Verbs      []string
	Fn         ActionFunc
	AuthRoles  []auth.Role
	name       string
	vanityPath string
}

func (a *Action) HasVerb(verb string) bool {
	return slices.Index(a.Verbs, verb) != -1
}

func (a *Action) GetFn() ActionFunc {
	return a.Fn
}

func (a *Action) AuthorizedRoles() []auth.Role {
	return a.AuthRoles
}

func (a *Action) VanityPath() string {
	return a.vanityPath
}

func (a *Action) WithAuthRoles(roles ...auth.Role) IAction {
	a.AuthRoles = roles

	return a
}

func (a *Action) WithVanityPath(path string) IAction {
	a.vanityPath = path

	return a
}

func (a *Action) Name() string {
	return a.name
}

func NewAction(name string, fn ActionFunc, verbs ...string) *Action {
	if len(verbs) < 1 {
		verbs = []string{http.MethodGet}
	}

	a := &Action{
		name:  name,
		Verbs: verbs,
		Fn:    fn,
		// default a new action to allow anyone to access it
		// requires authorization to be specified explicitly
		AuthRoles:  []auth.Role{auth.RoleAnonymous},
		vanityPath: "",
	}

	return a
}
