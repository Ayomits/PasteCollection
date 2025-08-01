package responses

type Violation struct {
	Message      string `json:"message"`
	PropertyPath string `json:"propertyPath"`
}

type ValidationError struct {
	Message    string      `json:"message"`
	Violations []Violation `json:"violations"`
}

func NewValidationError(message string, violations []Violation) *ValidationError {
	return &ValidationError{
		Message:    message,
		Violations: violations,
	}
}

func NewViolation(message string, propertyPath string) *Violation {
	return &Violation{
		Message:      message,
		PropertyPath: propertyPath,
	}
}
