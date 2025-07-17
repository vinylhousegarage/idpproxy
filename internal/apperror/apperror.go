package apperror

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) StatusCode() int {
	return e.Code
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
