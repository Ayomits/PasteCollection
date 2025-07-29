package responses

type UnauthorizedError struct {
	Error string `json:error,default="Unauthorized"`
}

func NewUnauthorizedError(err ...string) *UnauthorizedError {
    if len(err) > 0 {
        return &UnauthorizedError{
            Error: err[0],
        }
    }
    return &UnauthorizedError{
        Error: "Unauthorized",
    }
}
