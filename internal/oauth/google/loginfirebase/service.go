package loginfirebase

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}

var ErrInvalidGoogleLoginRequest = errors.New("invalid or missing google id_token")

func ParseGoogleLoginRequest(r *http.Request) (*GoogleLoginRequest, error) {
	if r.Body == nil {
		return nil, ErrInvalidGoogleLoginRequest

	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidGoogleLoginRequest
	}

	if req.IDToken == "" {
		return nil, ErrInvalidGoogleLoginRequest
	}

	return &req, nil
}
