package callback

import "testing"

func TestGitHubCallbackHandler_ready(t *testing.T) {
	t.Parallel()

	t.Run("returns_true_when_all_dependencies_are_set", func(t *testing.T) {
		t.Parallel()

		h := newHandlerForTest(
			t,
			&fakeHTTPClient{},
			&fakeUserService{},
			&fakeProxyCodeService{},
			&fakeGitHubTokenRepo{},
			&fakeSessionService{},
		)

		if !h.ready() {
			t.Fatal("ready() = false, want true")
		}
	})

	t.Run("returns_false_when_token_repository_is_nil", func(t *testing.T) {
		t.Parallel()

		h := newHandlerForTest(
			t,
			&fakeHTTPClient{},
			&fakeUserService{},
			&fakeProxyCodeService{},
			&fakeGitHubTokenRepo{},
			&fakeSessionService{},
		)
		h.TokenRepo = nil

		if h.ready() {
			t.Fatal("ready() = true, want false")
		}
	})

	t.Run("returns_false_when_session_service_is_nil", func(t *testing.T) {
		t.Parallel()

		h := newHandlerForTest(
			t,
			&fakeHTTPClient{},
			&fakeUserService{},
			&fakeProxyCodeService{},
			&fakeGitHubTokenRepo{},
			&fakeSessionService{},
		)
		h.SessionService = nil

		if h.ready() {
			t.Fatal("ready() = true, want false")
		}
	})
}
