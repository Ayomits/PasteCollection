package responses

type InternalError struct {
	Error string `json:error,default="Internal server exception"`
}

func NewInternalError(err ...string) *ForbiddenError {
	if len(err) > 0 {
		return &ForbiddenError{
			Error: err[0],
		}
	}
	return &ForbiddenError{
		Error: "Internal server exception",
	}
}
