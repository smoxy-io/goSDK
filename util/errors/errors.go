package errors

var (
	ErrNotFound = New("not found")
	ErrExists   = New("already exists")
	ErrInvalid  = New("invalid %s", "")
)
