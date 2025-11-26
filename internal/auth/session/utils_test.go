package session

import (
	"context"
	"errors"
)

type fakeRepository struct {
	created []*Session
	findMap map[string]*Session

	createErr error
	findErr   error
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		findMap: map[string]*Session{},
	}
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
	if f.findErr != nil {
		return nil, f.findErr
	}
	s := f.findMap[id]
	if s == nil {
		return nil, errors.New("not found")
	}

	return s, nil
}
