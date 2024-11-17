package auth

import "context"

const (
	IsAuthorizedContextKey = "isAuthorized"
	RoleContextKey         = "roles"
)

type Role int

const (
	RoleAnonymous Role = 0
	RoleUser      Role = 1 << iota
	RoleSupport
	RoleSales
	RoleDeveloper
	RoleAdmin
	RoleSuperAdmin
)

var (
	AllRoles       Role = RoleUser.Add(RoleSupport).Add(RoleSales).Add(RoleDeveloper).Add(RoleAdmin)
	SupportLevel   Role = RoleSupport.Add(RoleDeveloper).Add(RoleAdmin)
	SalesLevel     Role = RoleSales.Add(RoleDeveloper).Add(RoleAdmin)
	DeveloperLevel Role = RoleDeveloper.Add(RoleAdmin)
)

var (
	roles = []Role{
		RoleAnonymous,
		RoleUser,
		RoleSupport,
		RoleSales,
		RoleDeveloper,
		RoleAdmin,
		RoleSuperAdmin,
	}

	roleNameMap = map[Role]string{
		RoleAnonymous:  "anonymous",
		RoleUser:       "user",
		RoleSupport:    "support",
		RoleSales:      "sales",
		RoleDeveloper:  "developer",
		RoleAdmin:      "admin",
		RoleSuperAdmin: "superAdmin",
	}
	nameRoleMap = map[string]Role{
		"anonymous":  RoleAnonymous,
		"user":       RoleUser,
		"support":    RoleSupport,
		"sales":      RoleSales,
		"developer":  RoleDeveloper,
		"admin":      RoleAdmin,
		"superAdmin": RoleSuperAdmin,
	}

	authMap = map[string]Role{}
)

func (r Role) String() string {
	name, ok := roleNameMap[r]

	if !ok {
		if r > 0 {
			return "combined"
		}

		return "unknown"
	}

	return name
}

func (r Role) Has(role Role) bool {
	return r&role != 0
}

func (r Role) Add(role Role) Role {
	return r | role
}

func (r Role) Remove(role Role) Role {
	return r &^ role
}

func (r Role) ContextWithRoles(ctx context.Context) context.Context {
	return context.WithValue(ctx, RoleContextKey, r)
}

func NewRole(roles ...Role) Role {
	var r Role

	for _, role := range roles {
		r = r.Add(role)
	}

	return r
}

func NewRoleFromString(roles ...string) Role {
	r := []Role{}

	for _, strRole := range roles {
		if role, ok := nameRoleMap[strRole]; ok {
			r = append(r, role)
		}
	}

	return NewRole(r...)
}

func NewRoleFromApiKeyRole(roles ...ApiKeyRole) Role {
	r := []Role{}

	for _, protoRole := range roles {
		if role, ok := nameRoleMap[protoRole.GetName()]; ok {
			r = append(r, role)
		}
	}

	return NewRole(r...)
}

func IsAuthorized(route string, roles ...string) bool {
	role := NewRoleFromString(roles...)

	if role.Has(RoleSuperAdmin) {
		// super admin have access to everything
		return true
	}

	allowedRoles, ok := authMap[route]

	if !ok {
		return false
	}

	if allowedRoles == RoleAnonymous {
		// everyone is allowed
		return true
	}

	for _, r := range roles {
		if r2, rOk := nameRoleMap[r]; rOk {
			if allowedRoles.Has(r2) {
				return true
			}
		}
	}

	return false
}

func SetAllowedRoles(route string, roles ...Role) {
	r := NewRole(roles...)

	authMap[route] = r
}

func GetRolesFromCtx(ctx context.Context) Role {
	if role, ok := ctx.Value(RoleContextKey).(Role); ok {
		return role
	}

	return RoleAnonymous
}

func IsAuthorizedFromCtx(ctx context.Context) bool {
	isAuth, ok := ctx.Value(IsAuthorizedContextKey).(bool)

	if !ok {
		return false
	}

	return isAuth
}

func ContextWithIsAuthorized(ctx context.Context, route string) context.Context {
	if claims := GetUserClaimsFromCtx(ctx); claims != nil {
		return context.WithValue(ctx, IsAuthorizedContextKey, IsAuthorized(route, claims.Roles...))
	}

	// check if anonymous role is allowed
	return context.WithValue(ctx, IsAuthorizedContextKey, IsAuthorized(route, RoleAnonymous.String()))
}

func RegisterRole(name string) Role {
	if r, ok := nameRoleMap[name]; ok {
		// role is already registered
		return r
	}

	r := 1 << (len(roles) - 1)

	roles = append(roles, r)
	nameRoleMap[name] = r
	roleNameMap[r] = name

	AllRoles = AllRoles.Add(r)

	return r
}
