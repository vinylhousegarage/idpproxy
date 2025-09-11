package config

const (
	// for GitHub OAuth
	GitHubAllowSignup = "true"
	GitHubScope       = "read:user"

	// for GitHub API
	GitHubAPIBaseURL = "https://api.github.com"
	GitHubAPIVersion = "2022-11-28"
	GitHubUserURL    = GitHubAPIBaseURL + "/user"

	// for BuildGitHubLoginURL
	GitHubAuthorizeURL = "https://github.com/login/oauth/authorize"

	// for fetching Google OIDC metadata
	GoogleOIDCMetadataURL = "https://accounts.google.com/.well-known/openid-configuration"

	// for GenerateState
	StateLength = 16

	// for UserAgent
	UserAgentProduct = "idpproxy"
)
