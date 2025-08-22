package idp

type SignInWithIdpResp struct {
	ProviderID   string `json:"providerId,omitempty"`
	LocalID      string `json:"localId,omitempty"`
	IDToken      string `json:"idToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    string `json:"expiresIn,omitempty"`
	IsNewUser    bool   `json:"isNewUser,omitempty"`
}
