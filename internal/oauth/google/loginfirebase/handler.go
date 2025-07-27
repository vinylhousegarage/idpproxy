package loginfirebase

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/cookie"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/verify"
)

type LoginFirebaseHandler struct {
	Verifier verify.Verifier
	Logger   *zap.Logger
}

func NewLoginFirebaseHandler(
	verifier verify.Verifier,
	logger *zap.Logger,
) *LoginFirebaseHandler {
	return &LoginFirebaseHandler{
		Verifier: verifier,
		Logger:   logger,
	}
}

func (h *LoginFirebaseHandler) LoginFirebaseHandler(
	w http.ResponseWriter,
	r *http.Request,
) error {
	req, err := ParseGoogleLoginRequest(r)
	if err != nil {
		h.logger.Error("invalid request", zap.Error(err))

		return ErrInvalidRequest
	}

	_, err = verify.VerifyIDToken(r.Context(), h.Verifier, req.IDToken)
	if err != nil {
		h.logger.Error("unauthorized id_token", zap.Error(err))

		return ErrInvalidIDToken
	}

	cookie.SetIDTokenCookie(w, req.IDToken)
	w.WriteHeader(http.StatusOK)

	return nil
}
