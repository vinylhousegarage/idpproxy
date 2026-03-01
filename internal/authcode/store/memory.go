package store

import (
	"context"
	"sync"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type MemoryStore struct {
	mu         sync.Mutex
	proxyCodes map[string]authcode.ProxyCode
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		proxyCodes: make(map[string]authcode.ProxyCode),
	}
}

func (s *MemoryStore) Save(ctx context.Context, proxyCode authcode.ProxyCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.proxyCodes[proxyCode.Code] = proxyCode
	return nil
}

func (s *MemoryStore) Consume(ctx context.Context, proxyCodeValue, clientID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pc, ok := s.proxyCodes[proxyCodeValue]
	if !ok {
		return "", ErrNotFound
	}

	delete(s.proxyCodes, proxyCodeValue)

	if pc.ClientID != clientID {
		return "", ErrClientMismatch
	}

	if time.Now().After(pc.ExpiresAt) {
		return "", ErrExpired
	}

	return pc.UserID, nil
}
