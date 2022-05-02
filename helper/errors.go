package helper

const (
	ValidationError = iota + 1
	DatabaseError
	AuthenticationError
	DataDoesNotExistError
)

type CustomError struct {
	Code    int
	Message string
}

func (err *CustomError) Error() string {
	return err.Message
}

func NewCustomError(message string, code int) error {
	return &CustomError{
		Message: message,
		Code:    code,
	}
}
