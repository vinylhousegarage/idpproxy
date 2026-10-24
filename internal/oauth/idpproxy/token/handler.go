package token

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
	Logger  *zap.Logger
}

func NewHandler(svc *Service, logger *zap.Logger) *Handler {
	return &Handler{
		Service: svc,
		Logger:  logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn(
			"invalid token request",
			zap.Error(err),
		)

		http.Error(w, "invalid request", http.StatusBadRequest)

		return
	}

	resp, err := h.Service.Exchange(r.Context(), req)
	if err != nil {
		h.Logger.Warn(
			"token exchange failed",
			zap.String("client_id", req.ClientID),
			zap.String("grant_type", req.GrantType),
			zap.Error(err),
		)

		writeOAuthError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error(
			"encode token response failed",
			zap.Error(err),
		)
	}
}

func writeOAuthError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	type oauthErr struct {
		Error string `json:"error"`
	}

	status := http.StatusBadRequest
	var code string

	switch {
	case errors.Is(err, ErrUnsupportedGrantType):
		code = "unsupported_grant_type"
	case errors.Is(err, ErrInvalidClient):
		code = "invalid_client"
	case errors.Is(err, ErrInvalidGrant):
		code = "invalid_grant"
	default:
		status = http.StatusInternalServerError
		code = "server_error"
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(oauthErr{Error: code})
}
