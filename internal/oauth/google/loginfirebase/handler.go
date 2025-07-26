package loginfirebase

import (
	"net/http"

	"go.uber.org/zap"

	"github/vinylhousegarage/idpproxy/internal/cookie"
	"github/vinylhousegarage/idpproxy/internal/oauth/google/verify"
)

type LoginFirebaseHandler struct {
	Verifier verify.Verifier
}

func NewLoginFirebaseHandler(verifier verify.Verifier) *LoginFirebaseHandler {
	return &LoginFirebaseHandler{Verifier: verifier}
}

func (h *LoginFirebaseHandler) LoginFirebaseHandler(w http.ResponseWriter, r *http.Request, logger *zap.Logger) error {
	req, err := ParseGoogleLoginRequest(r)
	if err != nil {
		logger.Error("invalid request", zap.Error(err))

		return ErrInvalidRequest
	}

	_, err = verify.VerifyIDToken(r.Context(), h.Verifier, req.IDToken)
	if err != nil {
		logger.Error("unauthorized id_token", zap.Error(err))

		return ErrInvalidIDToken
	}

	cookie.SetIDTokenCookie(w, req.IDToken)
	w.WriteHeader(http.StatusOK)

	return nil
}
