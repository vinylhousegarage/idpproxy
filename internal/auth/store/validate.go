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

	provider, raw, ok := strings.Cut(userID, ":")
	if !ok {
		return fmt.Errorf("%w: want '<provider>:<uuid>'", ErrInvalidUserID)
	}
	provider = strings.ToLower(strings.TrimSpace(provider))

	if _, ok := allowedProviders[provider]; !ok {
		return fmt.Errorf("%w: unknown provider %q", ErrInvalidUserID, provider)
	}

	if _, err := uuid.Parse(strings.TrimSpace(raw)); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidUserID, err)
	}

	return nil
}

func validateFamilyID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("empty familyID")
	}
	if len(id) > 512 {
		return fmt.Errorf("familyID too long")
	}
	if containsSlash(id) {
		return fmt.Errorf("familyID must not contain '/'")
	}

	return nil
}

func containsSlash(s string) bool {
	return strings.Contains(s, "/")
}
