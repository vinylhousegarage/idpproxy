package me

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/httperror"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/verify"
)

type MeHandler struct {
	Verifier verify.Verifier
	Logger   *zap.Logger
}

func NewMeHandler(
	verifier verify.Verifier,
	logger *zap.Logger,
) *MeHandler {
	return &MeHandler{
		Verifier: verifier,
		Logger:   logger,
	}
}

func (h *MeHandler) Serve(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Writer.Header().Set("Vary", "Origin")
	c.Writer.Header().Set("Allow", "GET, OPTIONS")

	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusNoContent)
		return
	}

	idToken, err := ExtractAuthHeaderToken(c.Request)
	if err != nil {
		httperror.WriteErrorResponse(c.Writer, err, h.Logger)
		return
	}

	token, err := verify.VerifyIDToken(c.Request.Context(), h.Verifier, idToken)
	if err != nil {
		httperror.WriteErrorResponse(c.Writer, err, h.Logger)
		return
	}

	resp := MeResponse{
		Sub: token.UID,
		Iss: token.Claims["iss"].(string),
		Aud: token.Claims["aud"].(string),
		Exp: int64(token.Claims["exp"].(float64)),
	}

	c.JSON(http.StatusOK, resp)
}
