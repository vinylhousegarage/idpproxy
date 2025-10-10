package store

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
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

var allowedProviders = map[string]struct{}{
	"github": {},
	"google": {},
}

func validateUserID(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return fmt.Errorf("%w: empty", ErrInvalidUserID)
	}

	parts := strings.SplitN(userID, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: want '<provider>:<uuid>'", ErrInvalidUserID)
	}
	provider := parts[0]
	raw := parts[1]

	if _, ok := allowedProviders[provider]; !ok {
		return fmt.Errorf("%w: unknown provider %q", ErrInvalidUserID, provider)
	}

	if _, err := uuid.Parse(raw); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidUserID, err)
	}

	return nil
}
