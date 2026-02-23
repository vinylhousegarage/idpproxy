package callback

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/config"
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
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`not found`)),
		}, nil
	}
}

func okJSON(s string) *http.Response {
	h := make(http.Header)
	h.Set("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(s)),
		Header:     h,
	}
}

type fakeUserService struct {
	returnID string
	err      error
}

func (s *fakeUserService) UpsertFromGitHub(_ context.Context, _ int64, _ string, _ string) (string, error) {
	if s.err != nil {
		return "", s.err
	}

	return s.returnID, nil
}

type fakeAuthCodeService struct {
	code   string
	err    error
	called bool
}

func (f *fakeAuthCodeService) Issue(
	_ context.Context,
	_ string,
	_ string,
) (string, error) {
	f.called = true

	if f.err != nil {
		return "", f.err
	}

	return f.code, nil
}
