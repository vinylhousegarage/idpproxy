package store

import (
	"sync"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type MemoryStore struct {
	mu    sync.Mutex
	codes map[string]authcode.AuthCode
}
