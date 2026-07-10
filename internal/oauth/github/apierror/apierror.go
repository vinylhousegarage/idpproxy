package apierror

type ErrorCode string

type Internal struct {
	Code   ErrorCode
	Status int
	Err    error
}

type APIError struct {
	Code       ErrorCode
	HTTPStatus int
	Err        error
	Internal   []Internal
}

func (e *APIError) Error() string {
	return string(e.Code)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, status int, err error, internal ...Internal) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: status,
		Err:        err,
		Internal:   internal,
	}
}

func NewDetail(code ErrorCode, status int, err error) *Internal {
	return &Internal{
		Code:   code,
		Status: status,
		Err:    err,
	}
}
