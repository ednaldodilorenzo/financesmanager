package util

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidateRequestPayload[T any](parser func(out interface{}) error) (*T, error) {
	var errs []*ErrorResponse
	payload := new(T)

	if err := parser(&payload); err != nil {
		errs = append(errs, &ErrorResponse{Value: err.Error()})
		return nil, NewValidationError(err.Error(), errs)
	}

	err := validate.Struct(payload)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errs = append(errs, &element)
		}
		return nil, NewValidationError(err.Error(), errs)
	}

	return payload, nil
}
