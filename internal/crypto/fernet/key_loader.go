package fernet

import (
	"fmt"

	"github.com/fernet/fernet-go"
)

func DecodeFernetKey(raw string) (*fernet.Key, error) {
	var key fernet.Key
	if err := key.Decode(raw); err != nil {
		return nil, fmt.Errorf("failed to decode fernet key: %w", err)
	}
	return &key, nil
}
