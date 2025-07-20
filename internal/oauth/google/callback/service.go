package callback

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/httclient"

	"go.uber.org/zap"
)

func ValidateCallbackRequest(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		return "", ErrMissingCode
	}

	if state == "" {
		return "", ErrMissingState
	}

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return "", ErrMissingStateCookie
	}

	if state != cookie.Value {
		return "", ErrInvalidState
	}

	return code, nil
}

type callbackMetadata struct {
	TokenEndpoint string `json:"token_endpoint"`
}

func GetCallbackURL(
	metadataURL string,
	client httpclient.HTTPClient,
	logger *zap.Logger,
) (string, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		logger.Error("failed to create request", zap.String("url", metadataURL), zap.Error(err))
		return "", ErrFailedToCreateRequest
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch metadata", zap.String("url", metadataURL), zap.Error(err))
		return "", ErrFailedToFetchMetadata
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close response body", zap.Error(cerr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		fields := []zap.Field{
			zap.Int("status_code", resp.StatusCode),
			zap.String("url", req.URL.String()),
		}
		if readErr != nil {
			fields = append(fields, zap.Error(readErr))
			logger.Warn("failed to read response body", fields...)
		} else {
			fields = append(fields, zap.ByteString("body", body))
			logger.Warn("unexpected response from metadata", fields...)
		}
		return "", ErrUnexpectedMetadataStatusCode
	}

	var meta callbackMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", ErrFailedToDecodeMetadata
	}

	if meta.TokenEndpoint == "" {
		logger.Error("missing token_endpoint")
		return "", ErrMissingTokenEndpoint
	}

	logger.Info("token_endpoint retrieved successfully", zap.String("token_endpoint", meta.TokenEndpoint))
	return meta.TokenEndpoint, nil
}
