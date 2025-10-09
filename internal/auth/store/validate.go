package store

import (
	"fmt"
	"strings"
)

func validateRefreshID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: empty", ErrInvalidID)
	}
	if strings.Contains(id, "/") {
		return fmt.Errorf("%w: must not contain '/'", ErrInvalidID)
	}

	return nil
}
