package logout

import (
	"context"
	"net/http"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

type SessionInvalidator interface {
	Invalidate(
		ctx context.Context,
		sessionID string,
	) (*authsession.Session, error)
}

type Handler struct {
	SessionInvalidator SessionInvalidator
}

func NewHandler(invalidator SessionInvalidator) *Handler {
	return &Handler{
		SessionInvalidator: invalidator,
	}
}

func (h *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	cookie, err := r.Cookie(authsession.SessionCookieName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(
			w,
			"invalid session cookie",
			http.StatusBadRequest,
		)
		return
	}

	if _, err := h.SessionInvalidator.Invalidate(
		r.Context(),
		cookie.Value,
	); err != nil {
		http.Error(
			w,
			"failed to logout",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authsession.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusNoContent)
}
