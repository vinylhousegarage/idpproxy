package httpclient

import (
	"fmt"
	"io"
	"net/http"
)

func ValidateResponse(resp *http.Response) (int, error) {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.StatusCode, nil
	}

	snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 256))

	return resp.StatusCode, fmt.Errorf("upstream error: %s", snippet)
}
