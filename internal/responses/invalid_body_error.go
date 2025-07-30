package responses

type InvalidBodyError struct {
	Message    string `json:"message"`
	violations any    `json:"violations"`
}

func NewInvalidBodyError(err any) *InvalidBodyError {
	return &InvalidBodyError{
		Message:    "Invalid body",
		violations: err,
	}
}
