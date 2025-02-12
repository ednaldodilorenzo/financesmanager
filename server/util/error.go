package util

import "fmt"

type ItemNotFoundError struct {
	Message string
}

func (e *ItemNotFoundError) Error() string {
	return e.Message
}

// CustomDatabaseError represents a custom error type for database issues.
type RuntimeError struct {
	Msg string
	Err error
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Database error: %s | Original error: %v", e.Msg, e.Err)
}

func (e *RuntimeError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new CustomDatabaseError.
func NewRuntimeError(msg string, err error) *RuntimeError {
	return &RuntimeError{
		Msg: msg,
		Err: err,
	}
}

const (
	BE_INPUT_VALIDATION_ERROR int = iota
	BE_PASSWORD_DO_NOT_MATCH
	BE_USER_ALREADY_REGISTERED
	BE_USER_EMAIL_NOT_VERIFIED
)

type BusinessError struct {
	Message string
	Err     error
	Code    int
}

func (b *BusinessError) Error() string {
	return fmt.Sprintf("Business error: %s, code %d", b.Message, b.Code)
}

func NewBusinessError(msg string, err error, code int) *BusinessError {
	return &BusinessError{
		Message: msg,
		Err:     err,
		Code:    code,
	}
}
