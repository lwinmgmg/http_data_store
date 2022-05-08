package helper

import (
	"fmt"
)

const (
	ValidationError = iota
	DatabaseError
	AuthenticationError
	DataDoesNotExistError
	InternalServerError
)

var HttpErrorCodes = [...]int{
	422,
	500,
	401,
	404,
	500,
}

func ErrorResponse(input any) (int, map[string]any) {
	if v, ok := input.(CustomError); ok {
		return HttpErrorCodes[v.Code], map[string]any{"detail": v.Message}
	}
	return 500, map[string]any{"detail": fmt.Sprintf("Internal Server Error : %v", input)}
}

type CustomError struct {
	Code    int
	Message string
}

func (err *CustomError) Error() string {
	return err.Message
}

func NewCustomError(code int, message string, args ...any) error {
	return &CustomError{
		Message: fmt.Sprintf(message, args...),
		Code:    code,
	}
}
