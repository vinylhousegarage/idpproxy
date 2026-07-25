package apierror

import "fmt"

func (e *APIError) AddInternal(code ErrorCode, key, value string) *APIError {
	if code == "" || value == "" {
		return e
	}

	e.Internals = append(e.Internals, APIInternal{
		Code:   code,
		Status: 500,
		Err:    fmt.Errorf("%s: %s", key, value),
	})

	return e
}

func (e *APIError) GetHTTPStatus() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	if len(e.Internals) > 0 {
		return e.Internals[0].Status
	}

	return 500
}
