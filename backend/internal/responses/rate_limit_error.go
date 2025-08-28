package responses

type RateLimitError struct {
	Message string `json:"message"`
}

func NewRateLimitError(err ...string) *RateLimitError {
	if len(err) > 0 {
		return &RateLimitError{
			Message: err[0],
		}
	}
	return &RateLimitError{
		Message: "You are being rate limit...",
	}
}
