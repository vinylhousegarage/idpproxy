package loginfirebase

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
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
		h.Logger.Error("invalid request", zap.Error(err))

		return ErrInvalidRequest
	}

	_, err = verify.VerifyIDToken(r.Context(), h.Verifier, req.IDToken)
	if err != nil {
		h.Logger.Error("unauthorized id_token", zap.Error(err))

		return ErrInvalidIDToken
	}

	cookie.SetIDTokenCookie(w, req.IDToken)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *LoginFirebaseHandler) Serve(c *gin.Context) {
	if err := h.LoginFirebaseHandler(c.Writer, c.Request); err != nil {
		h.Logger.Warn("loginfirebase failed", zap.Error(err))

		switch err {
		case ErrInvalidRequest:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
		case ErrInvalidIDToken:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized id_token",
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
	}
}
