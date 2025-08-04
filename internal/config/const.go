package config

const (
	// for fetching Google OIDC metadata
	GoogleOIDCMetadataURL = "https://accounts.google.com/.well-known/openid-configuration"

	// for GoogleConfig
	GoogleResponseType = "code"
	GoogleScope        = "openid email profile"
	GoogleAccessType   = "offline"
	GooglePrompt       = "consent"

	// for GoogleTokenRepository
	CollectionGoogleTokens = "google_tokens"

	// for GenerateState
	StateLength = 16

	// for BuildGitHubLoginURL
	GitHubAuthorizeURL = "https://github.com/login/oauth/authorize"
)
