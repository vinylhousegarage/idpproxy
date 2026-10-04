package callback

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/config"
	githubstore "github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
)

type fakeHTTPClient struct {
	tokenJSON     string
	userJSON      string
	forceTokenErr bool
	forceUserErr  bool
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	switch req.URL.String() {
	case config.GitHubTokenURL:
		if f.forceTokenErr {
			return nil, errors.New("token http error")
		}

		return okJSON(f.tokenJSON), nil

	case config.GitHubUserURL:
		if f.forceUserErr {
			return nil, errors.New("user http error")
		}

		return okJSON(f.userJSON), nil

	default:
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(`not found`)),
		}, nil
	}
}

func okJSON(s string) *http.Response {
	h := make(http.Header)
	h.Set("Content-Type", "application/json")

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(s)),
		Header:     h,
	}
}

type fakeUserService struct {
	returnID string
	err      error
}

func (s *fakeUserService) UpsertFromGitHub(
	_ context.Context,
	_ int64,
	_ string,
	_ string,
) (string, error) {
	if s.err != nil {
		return "", s.err
	}

	return s.returnID, nil
}

type fakeProxyCodeService struct {
	proxyCode string
	err       error
	called    bool
}

func (f *fakeProxyCodeService) Issue(
	_ context.Context,
	_ string,
	_ string,
) (string, error) {
	f.called = true

	if f.err != nil {
		return "", f.err
	}

	return f.proxyCode, nil
}

type fakeGitHubTokenRepo struct {
	called bool
	rec    *githubstore.GitHubTokenRecord
	err    error
}

func (r *fakeGitHubTokenRepo) Upsert(
	_ context.Context,
	rec *githubstore.GitHubTokenRecord,
) error {
	r.called = true
	r.rec = rec

	return r.err
}

type fakeSessionService struct {
	session *authsession.Session
	err     error

	called bool
	userID string
}

func (s *fakeSessionService) Start(
	_ context.Context,
	userID string,
) (*authsession.Session, error) {
	s.called = true
	s.userID = userID

	if s.err != nil {
		return nil, s.err
	}

	return s.session, nil
}
