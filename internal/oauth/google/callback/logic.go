package callback

import (
	"fmt"
	"net/url"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func BuildTokenRequestBody(code string) (string, error) {
	googleConfig, err := config.LoadGoogleConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load Google config: %w", err)
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", googleConfig.ClientID)
	form.Set("client_secret", googleConfig.ClientSecret)
	form.Set("redirect_uri", googleConfig.RedirectURI)

	return form.Encode(), nil
}
