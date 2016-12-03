package auth

import "time"

type AuthenticationToken struct {
	Tokentype    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreationTime time.Time
	Error        error
}

func (vSelf *AuthenticationToken) IsExpired() bool {

	return false
}

func (vSelf *AuthenticationToken) IsValid() bool {
	return vSelf.Error == nil
}
