package responses

type ForbiddenError struct {
	Error string `json:error,default="Internal server exception"`
}

func NewForbiddenError(err ...string) *ForbiddenError {
	if len(err) > 0 {
		return &ForbiddenError{
			Error: err[0],
		}
	}
	return &ForbiddenError{
		Error: "Access denied",
	}
}
