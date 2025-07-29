package me

import (
	"encoding/json"
	"net/http"

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

func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	idToken, err := ExtractAuthHeaderToken(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, h.Logger)
		return
	}

	token, err := verify.VerifyIDToken(r.Context(), h.Verifier, idToken)
	if err != nil {
		httperror.WriteErrorResponse(w, err, h.Logger)
		return
	}

	resp := MeResponse{
		Sub: token.UID,
		Iss: token.Claims["iss"].(string),
		Aud: token.Claims["aud"].(string),
		Exp: int64(token.Claims["exp"].(float64)),
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error("failed to write user response", zap.Error(err))
		httperror.WriteErrorResponse(w, ErrFailedToWriteUserResponse, h.Logger)
	}
}
