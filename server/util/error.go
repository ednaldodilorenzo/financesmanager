package util

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

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
	BE_NOT_FOUND
)

type Error struct {
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return e.Message
}

func NewError(statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}

var (
	ErrNotFound     = NewError(fiber.StatusNotFound, "item not found")
	ErrBadRequest   = NewError(fiber.StatusBadRequest, "bad request")
	ErrBusiness     = NewError(fiber.StatusUnprocessableEntity, "business error")
	ErrUnauthorized = NewError(fiber.StatusUnauthorized, "unauthorized")
)

type APIError struct {
	Messages []string `json:"messages"`
	Err      *Error
}

func (ae *APIError) Error() string {
	return ae.Err.Error()
}

func NewAPIError(err *Error, messages []string) *APIError {
	return &APIError{
		Err:      err,
		Messages: messages,
	}
}

func ServerErrorHandler(ctx *fiber.Ctx, err error) error {

	var apiError *APIError
	if errors.As(err, &apiError) {
		return ctx.Status(apiError.Err.StatusCode).JSON(fiber.Map{"status": "fail", "errors": apiError.Messages})
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		return ctx.Status(fiberError.Code).JSON(fiber.Map{"status": "fail", "errors": fiberError.Message})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": "Internal Server Error"})
}
