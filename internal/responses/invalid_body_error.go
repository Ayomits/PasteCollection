package responses

type InvalidBodyError struct {
	Error      string `json:error`
	violations any    `json:violations`
}

func NewInvalidBodyError(err any) *InvalidBodyError {
	return &InvalidBodyError{
		Error: "Invalid body",
		violations: err,
	}
}
