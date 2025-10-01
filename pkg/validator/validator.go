package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		return v.formatError(err)
	}
	return nil
}

func (v *Validator) formatError(err error) error {
	var messages []string

	for _, err := range err.(validator.ValidationErrors) {
		var message string
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "email":
			message = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters", field, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}

		messages = append(messages, message)
	}

	return fmt.Errorf("%s", strings.Join(messages, ", "))
}