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

type NotFoundError struct {
	Message string
}

func (n *NotFoundError) Error() string {
	return "Not found"
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		Message: msg,
	}
}

type ValidationError struct {
	Message string
	Errors  []*ErrorResponse
}

func (n *ValidationError) Error() string {
	return n.Message
}

func NewValidationError(message string, errors []*ErrorResponse) *ValidationError {
	return &ValidationError{
		Message: message,
		Errors:  errors,
	}
}

type UnauthorizedError struct {
	Message string
}

func (n *UnauthorizedError) Error() string {
	return n.Message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: message,
	}
}

func ServerErrorHandler(ctx *fiber.Ctx, err error) error {
	var notFoundError *NotFoundError
	if errors.As(err, &notFoundError) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "errors": notFoundError.Message})
	}

	var businessError *BusinessError
	if errors.As(err, &businessError) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": businessError.Message, "code": businessError.Code})
	}

	var validationError *ValidationError
	if errors.As(err, &validationError) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": validationError.Errors})
	}

	var unauthorizedError *UnauthorizedError
	if errors.As(err, &unauthorizedError) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "errors": unauthorizedError.Error()})
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		return ctx.Status(fiberError.Code).JSON(fiber.Map{"status": "fail", "errors": fiberError.Message})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "errors": "Internal Server Error"})
}
