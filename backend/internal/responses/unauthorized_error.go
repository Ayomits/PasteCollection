package responses

type UnauthorizedError struct {
	Message string `json:"message"`
}

func NewUnauthorizedError(err ...string) *UnauthorizedError {
    if len(err) > 0 {
        return &UnauthorizedError{
            Message: err[0],
        }
    }
    return &UnauthorizedError{
        Message: "Unauthorized",
    }
}
