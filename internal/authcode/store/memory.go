package store

import (
	"context"
	"sync"
	"time"

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

func (s *MemoryStore) Consume(ctx context.Context, code, clientID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ac, ok := s.codes[code]
	if !ok {
		return "", ErrNotFound
	}

	delete(s.codes, code)

	if ac.ClientID != clientID {
		return "", ErrClientMismatch
	}

	if time.Now().After(ac.ExpiresAt) {
		return "", ErrExpired
	}

	return ac.UserID, nil
}
