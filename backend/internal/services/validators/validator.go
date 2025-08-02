package validators

import (
	"api/internal/responses"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AppValidator interface {
	Validate(body any) *responses.ValidationError
}

type validationFormatter struct {
	validate *validator.Validate
}

func NewAppValidator() AppValidator {
	return &validationFormatter{
		validate: validator.New(),
	}
}

func (f *validationFormatter) Validate(body any) *responses.ValidationError {
	err := f.validate.Struct(body)
	var violations []responses.Violation

	if err == nil {
		return nil
	}

	for _, violation := range err.(validator.ValidationErrors) {
		violations = append(violations, *responses.NewViolation(f.formatErrorMessage(violation), f.getJSONFieldName(violation, body)))
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

func (f *validationFormatter) formatErrorMessage(fieldErr validator.FieldError) string {
	return fieldErr.Error()
}

var AppValidatorInstance AppValidator = NewAppValidator()
