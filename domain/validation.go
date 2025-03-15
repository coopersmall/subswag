package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

func validatePermission(fl validator.FieldLevel) bool {
	permission := Permission(fl.Field().String())
	return permission.Validate() == nil
}

func validateNotZeroTime(fl validator.FieldLevel) bool {
	timeVal, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return !timeVal.IsZero()
}

func Validate(i any) error {
	v := validator.New()
	err := v.RegisterValidation("permission", validatePermission)
	if err != nil {
		return err
	}
	err = v.RegisterValidation("nonZeroTime", validateNotZeroTime)
	if err != nil {
		return err
	}
	err = v.Struct(i)
	if err == nil {
		return nil
	}

	fieldValidationErrors := []FieldValidationError{}
	for _, e := range err.(validator.ValidationErrors) {
		fieldValidationErrors = append(fieldValidationErrors, FieldValidationError{
			Field: e.Field(),
			Error: e.Tag(),
		})
	}
	return NewValidationError(fieldValidationErrors...)
}

type ValidationError struct {
	Fields []FieldValidationError
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid fields: %v", e.Fields)
}

func NewValidationError(fields ...FieldValidationError) *ValidationError {
	return &ValidationError{Fields: fields}
}

type FieldValidationError struct {
	Field string
	Error string
}
