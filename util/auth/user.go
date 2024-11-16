package auth

type User interface {
	GetId() string
	GetStringRoles() []string
}
