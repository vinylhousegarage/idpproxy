package token

import "github.com/google/uuid"

func newFamilyID() string {
	return uuid.NewString()
}
