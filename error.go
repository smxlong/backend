package backend

// Error of the backend.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrTypeAlreadyRegistered is returned when a type is already registered,
	// such as when adding a singleton.
	ErrTypeAlreadyRegistered = Error("type already registered")
	// ErrTypeNotRegistered is returned when a type is not registered, such as
	// when getting a singleton.
	ErrTypeNotRegistered = Error("type not registered")
	// ErrInvalidFactory is returned when a factory is invalid, such as when
	// adding a transient.
	ErrInvalidFactory = Error("invalid factory")
	// ErrNotInvokable is returned when a function passed to Invoke is not
	// invokable.
	ErrNotInvokable = Error("not invokable")
)
