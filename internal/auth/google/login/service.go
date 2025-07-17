package login

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/httpclient"

	"go.uber.org/zap"
)

type GoogleLoginMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

func GetGoogleLoginURL(
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
		body, _ := io.ReadAll(resp.Body)
		logger.Error("unexpected response from metadata",
			zap.Int("status", resp.StatusCode),
			zap.ByteString("body", body),
		)
		return "", ErrUnexpectedMetadataStatusCode
	}

	var meta GoogleLoginMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", ErrFailedToDecodeMetadata
	}

	if meta.AuthorizationEndpoint == "" {
		logger.Error("missing authorization_endpoint")
		return "", ErrMissingAuthorizationEndpoint
	}

	logger.Info("authorization_endpoint retrieved successfully", zap.String("authorization_endpoint", meta.AuthorizationEndpoint))
	return meta.AuthorizationEndpoint, nil
}
