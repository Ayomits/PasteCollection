package responses

type NotFoundError struct {
	Message string `json:"message"`
}

func NewNotFoundError(err ...string) *NotFoundError {
	if len(err) > 0 {
		return &NotFoundError{
			Message: err[0],
		}
	}
	return &NotFoundError{
		Message: "Not found",
	}
}
