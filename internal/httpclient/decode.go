package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSON[T any](resp *http.Response, target *T) error {
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
