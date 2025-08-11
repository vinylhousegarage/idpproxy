package response

type AuthResultResponse struct {
    Sub      string `json:"sub"`
    Provider string `json:"provider"`
    Status   int    `json:"status"`
}

type UserProfileResponse struct {
    Sub      string `json:"sub"`
    Provider string `json:"provider"`
    Status   int    `json:"status"`
    Login    string `json:"login,omitempty"`
    Name     string `json:"name,omitempty"`
    Email    string `json:"email,omitempty"`
}

type GitHubUserAPIResponse struct {
    ID    int64  `json:"id"`
    Login string `json:"login"`
    Email string `json:"email,omitempty"`
    Name  string `json:"name,omitempty"`
}
