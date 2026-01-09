package store

import (
	"context"
	"sync"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type MemoryStore struct {
	mu    sync.Mutex
	codes map[string]authcode.AuthCode
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		codes: make(map[string]authcode.AuthCode),
	}
}

func (s *MemoryStore) Save(ctx context.Context, code authcode.AuthCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.codes[code.Code] = code
	return nil
}

func (s *MemoryStore) Get(
	ctx context.Context,
	code string,
) (authcode.AuthCode, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ac, ok := s.codes[code]
	if !ok {
		return authcode.AuthCode{}, ErrNotFound
	}

	return ac, nil
}
