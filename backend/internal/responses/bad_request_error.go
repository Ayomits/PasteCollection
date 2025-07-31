package responses

type BadRequestError struct {
	Message string `json:"message"`
}

func NewBadRequestError(err ...string) *BadRequestError {
	if len(err) > 0 {
		return &BadRequestError{
			Message: err[0],
		}
	}
	return &BadRequestError{
		Message: "Bad Request",
	}
}
