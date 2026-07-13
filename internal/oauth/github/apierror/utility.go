package apierror

import (
	"errors"
	"fmt"
)

func FormatDetail(key string, value string) (string, error) {
	if key == "" || value == "" {
		return "", errors.New("key and value must not be empty")
	}

	return fmt.Sprintf("%s: %s", key, value), nil
}

func (e *APIError) GetHTTPStatus() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	if len(e.Internal) > 0 {
		return e.Internal[0].Status
	}

	return 500
}
