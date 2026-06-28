package apierror

type ErrorCode string

type APIError struct {
	Code       ErrorCode
	HTTPStatus int
	Err        error
	Internal   []string
}

func (e *APIError) Error() string {
	return string(e.Code)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, status int, err error, internal ...string) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: status,
		Err:        err,
		Internal:   internal,
	}
}
