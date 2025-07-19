package fernet

import (
	"fmt"

	"github.com/fernet/fernet-go"
)

func DecodeFernetKey(raw string) (*fernet.Key, error) {
	keys, err := fernet.DecodeKeys(raw)
	if err != nil || len(keys) == 0 || keys[0] == nil {
		return nil, fmt.Errorf("failed to decode fernet key: %w", err)
	}

	return keys[0], nil
}
