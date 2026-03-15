package apierror

type ErrorCode string

type APIError struct {
	Code       ErrorCode
	HTTPStatus int
	Err        error
}

func (e *APIError) Error() string {
	return string(e.Code)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, status int, err error) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: status,
		Err:        err,
	}
}
