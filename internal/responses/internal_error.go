package responses

type InternalError struct {
	Message string `json:"message"`
}

func NewInternalError(err ...string) *ForbiddenError {
	if len(err) > 0 {
		return &ForbiddenError{
			Message: err[0],
		}
	}
	return &ForbiddenError{
		Message: "Internal server exception",
	}
}
