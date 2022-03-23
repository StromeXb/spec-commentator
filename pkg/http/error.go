package http

type Error struct {
	message string
}

func NewError(message string) error {
	return &Error{
		message: message,
	}
}

func (e *Error) Error() string {
	return e.message
}
