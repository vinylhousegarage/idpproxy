package session

import (
	"context"
	"errors"
	"time"
)

type Usecase struct {
	Repo        Repository
	Now         func() time.Time
	TTL         time.Duration
	IDGenerator func() (string, error)
}

func (uc *Usecase) Start(ctx context.Context, userID string) (*Session, error) {
	if uc == nil || uc.Repo == nil || uc.Now == nil || uc.IDGenerator == nil || uc.TTL <= 0 {
		return nil, errors.New("session: invalid usecase configuration")
	}
	if userID == "" {
		return nil, errors.New("session: empty userID")
	}

	now := uc.Now().UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	expiresAt := now.Add(uc.TTL)

	sessionID, err := uc.IDGenerator()
	if err != nil {
		return nil, err
	}

	s := &Session{
		SessionID: sessionID,
		UserID:    userID,
		Status:    "active",
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	if err := uc.Repo.Create(ctx, s); err != nil {
		return nil, err
	}

	return s, nil
}
