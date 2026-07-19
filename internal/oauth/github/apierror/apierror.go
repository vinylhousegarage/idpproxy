package apierror

type ErrorCode string

type APIInternal struct {
	Code   ErrorCode
	Status int
	Err    error
}

type APIError struct {
	Code       ErrorCode
	HTTPStatus int
	Err        error
	Internals  []APIInternal
}

func (e *APIError) Error() string {
	return string(e.Code)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, status int, err error, internals ...APIInternal) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: status,
		Err:        err,
		Internals:  internals,
	}
}

func NewAPIInternal(code ErrorCode, status int, err error) *APIInternal {
	return &APIInternal{
		Code:   code,
		Status: status,
		Err:    err,
	}
}
