package responses

type ForbiddenError struct {
	Message string `json:"message"`
}

func NewForbiddenError(err ...string) *ForbiddenError {
	if len(err) > 0 {
		return &ForbiddenError{
			Message: err[0],
		}
	}
	return &ForbiddenError{
		Message: "Access denied",
	}
}
