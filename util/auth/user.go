package auth

type User interface {
	GetId() string
	GetStringRoles() []string
}

type Scope interface {
	GetScope() string
}

type Extra interface {
	GetExtra() map[string]string
}
