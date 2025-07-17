package config

const GoogleOIDCMetadataURL = "https://accounts.google.com/.well-known/openid-configuration"

const (
	GoogleResponseType = "code"
	GoogleScope        = "openid email profile"
	GoogleAccessType   = "offline"
	GooglePrompt       = "consent"
)
