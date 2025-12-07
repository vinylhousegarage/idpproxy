package session

import (
	"context"
	"time"
)

type fakeRepository struct {
	created []*Session
	findMap map[string]*Session
	updated []*Session

	createErr error
	findErr   error
	updateErr error

	lastFindID string

	purgeErr    error
	purgeBefore time.Time
	purgeResult int
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		findMap: map[string]*Session{},
	}
}

func newFakeRepositoryWithSession(s *Session) *fakeRepository {
	r := newFakeRepository()
	if s != nil {
		r.findMap[s.SessionID] = s
	}
	return r
}

func (f *fakeRepository) Create(_ context.Context, s *Session) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, s)
	f.findMap[s.SessionID] = s

	return nil
}

func (f *fakeRepository) FindByID(_ context.Context, id string) (*Session, error) {
	f.lastFindID = id

	if f.findErr != nil {
		return nil, f.findErr
	}
	s := f.findMap[id]
	if s == nil {
		return nil, ErrNotFound
	}

	return s, nil
}

func (f *fakeRepository) Update(_ context.Context, s *Session) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.updated = append(f.updated, s)
	f.findMap[s.SessionID] = s

	return nil
}

func (f *fakeRepository) PurgeExpired(_ context.Context, before time.Time) (int, error) {
	if f.purgeErr != nil {
		return 0, f.purgeErr
	}
	f.purgeBefore = before

	return f.purgeResult, nil
}
