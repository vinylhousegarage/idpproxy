package loginfirebase

import (
	"net/http"

	"go.uber.org/zap"

	"github/vinylhousegarage/idpproxy/internal/cookie"
	"github/vinylhousegarage/idpproxy/internal/oauth/google/verify"
)

func (h *Handler) LoginFirebaseHandler(w http.ResponseWriter, r *http.Request, logger *zap.Logger) error {
	req, err := ParseIDTokenRequest(r)
	if err != nil {
		logger.Error("invalid request", zap.Error(err))

		return ErrInvalidRequest
	}

	_, err := verify.VerifyIDToken(r.Context(), h.verifier, req.IDToken)
	if err != nil {
		logger.Error("unauthorized id_token", zap.Error(err))

		return ErrInvalidIDToken
	}

	cookie.SetIDTokenCookie(w, req.IDToken)
	w.WriteHeader(http.StatusOK)

	return nil
}
