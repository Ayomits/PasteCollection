package formatter

import (
	"api/internal/responses"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationFormatter interface {
	Validate(body any) *responses.ValidationError
}

type validationFormatter struct {
	validate *validator.Validate
}

func NewValidationFormatter(v *validator.Validate) ValidationFormatter {
	return &validationFormatter{
		validate: v,
	}
}

func (f *validationFormatter) Validate(body any) *responses.ValidationError {
	err := f.validate.Struct(body)
	var violations []responses.Violation

	if err == nil {
		return nil
	}

	for _, violation := range err.(validator.ValidationErrors) {
		violations = append(violations, *responses.NewViolation(f.formatErrorMessage(violation, body), f.getJSONFieldName(violation, body)))
	}

	return responses.NewValidationError("Invalid payload", violations)
}

func (f *validationFormatter) getJSONFieldName(fieldErr validator.FieldError, body any) string {
	t := reflect.TypeOf(body)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldErr.Field())
	if !found {
		return fieldErr.Field()
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldErr.Field()
	}

	return strings.Split(jsonTag, ",")[0]
}

func (f *validationFormatter) formatErrorMessage(fieldErr validator.FieldError, body any) string {
	fiedName := f.getJSONFieldName(fieldErr, body)
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", fiedName)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fiedName, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be no more than %s characters", fiedName, fieldErr.Param())
	case "email":
		return "Must be a valid email address"
	case "oneof":
		options := strings.ReplaceAll(fieldErr.Param(), " ", ", ")
		return fmt.Sprintf("%s must be one of: %s", fiedName, options)
	default:
		return fmt.Sprintf("Validation failed for field %s", fiedName)
	}
}
