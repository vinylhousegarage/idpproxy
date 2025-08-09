package firebase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

var (
    ErrEmptyBody          = errors.New("empty response body")
    ErrMissingAccessToken = errors.New("missing access_token in response")
    ErrGitHubOAuthError   = errors.New("github oauth error")
    ErrNon2xx             = errors.New("non-2xx response")
    ErrParseFormBody      = errors.New("parse form body")
)

const (
    maxReadBytes int64 = 1 << 20
    snippetLimit       = 4096
)

type accessTokenJSON struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

func ExtractAccessTokenFromResponse(resp *http.Response) (string, error) {
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        snippet, _ := io.ReadAll(io.LimitReader(resp.Body, snippetLimit))
        return "", fmt.Errorf("%w: status=%d body=%q", ErrNon2xx, resp.StatusCode, string(snippet))
    }

    body, err := io.ReadAll(io.LimitReader(resp.Body, maxReadBytes))
    if err != nil {
        return "", fmt.Errorf("read response: %w", err)
    }
    if len(body) == 0 {
        return "", ErrEmptyBody
    }

    contentType := strings.TrimSpace(strings.ToLower(resp.Header.Get("content-type")))
    mediaType, _, _ := mime.ParseMediaType(contentType)

    switch {
    case mediaType == "application/json" || strings.HasSuffix(mediaType, "+json"):
        var accessTokenResp accessTokenJSON
        if err := json.Unmarshal(body, &accessTokenResp); err != nil {
            return "", fmt.Errorf("decode json: %w", err)
        }
        if accessTokenResp.Error != "" {
            return "", fmt.Errorf("%w: %s (%s)", ErrGitHubOAuthError, accessTokenResp.Error, accessTokenResp.ErrorDescription)
        }
        token := strings.TrimSpace(accessTokenResp.AccessToken)
        if token == "" {
            return "", ErrMissingAccessToken
        }
        return token, nil

    default:
        values, err := url.ParseQuery(string(body))
        if err != nil {
            return "", fmt.Errorf("%w: %v", ErrParseFormBody, err)
        }
        if e := strings.TrimSpace(values.Get("error")); e != "" {
            return "", fmt.Errorf("%w: %s (%s)", ErrGitHubOAuthError, e, values.Get("error_description"))
        }
        token := strings.TrimSpace(values.Get("access_token"))
        if token == "" {
            return "", ErrMissingAccessToken
        }
        return token, nil
    }
}
