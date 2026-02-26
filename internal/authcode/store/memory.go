package store

import (
	"context"
	"sync"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type MemoryStore struct {
	mu        sync.Mutex
	authCodes map[string]authcode.AuthCode
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		authCodes: make(map[string]authcode.AuthCode),
	}
}

func (s *MemoryStore) Save(ctx context.Context, authCode authcode.AuthCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.authCodes[authCode.Code] = authCode
	return nil
}

func (s *MemoryStore) Consume(ctx context.Context, authCodeValue, clientID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	authCode, ok := s.authCodes[authCodeValue]
	if !ok {
		return "", ErrNotFound
	}

	delete(s.authCodes, authCodeValue)

	if authCode.ClientID != clientID {
		return "", ErrClientMismatch
	}

	if time.Now().After(authCode.ExpiresAt) {
		return "", ErrExpired
	}

	return authCode.UserID, nil
}
