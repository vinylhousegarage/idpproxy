package idp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

const providerGitHub = "github.com"

func SignInGitHubWithAccessToken(
	ctx context.Context,
	httpClient httpclient.HTTPClient,
	apiKey string,
	requestURI string,
	accessToken string,
) (*signInGitHubWithAccessTokenResp, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey is empty")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("accessToken is empty")
	}
	if requestURI == "" {
		return nil, fmt.Errorf("requestURI is empty (must be an authorized domain)")
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	endpoint := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithIdp?key=" +
		url.QueryEscape(apiKey)

	pb := url.Values{}
	pb.Set("access_token", accessToken)
	pb.Set("providerId", providerGitHub)

	payload := signInPayload{
		RequestURI:        requestURI,
		PostBody:          pb.Encode(),
		ReturnSecureToken: true,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(io.LimitReader(resp.Body, 4<<10)).Decode(&e)
		if e.Error.Message != "" {
			return nil, fmt.Errorf("signInWithIdp: %s (status=%d)", e.Error.Message, resp.StatusCode)
		}
		return nil, fmt.Errorf("signInWithIdp: unexpected status %d", resp.StatusCode)
	}

	var out signInGitHubWithAccessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &out, nil
}
