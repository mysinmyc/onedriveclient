package auth

import (
	"fmt"
	"time"
)

type AuthenticationToken struct {
	Tokentype    string        `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
	Scope        string        `json:"scope"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	CreationTime time.Time
	Error        error
}

type TokenExpiredError struct {
	error
}

func (vSelf *TokenExpiredError) Error() string {
	return "Authentication token expired"
}

func newAuthenticationToken() *AuthenticationToken {
	return &AuthenticationToken{CreationTime: time.Now()}
}

func newAuthenticationTokenError(pError error) *AuthenticationToken {
	return &AuthenticationToken{CreationTime: time.Now(), Error: pError}
}

func (vSelf *AuthenticationToken) IsExpired() bool {
	return time.Now().After(vSelf.CreationTime.Add(time.Second * vSelf.ExpiresIn))
}

func (vSelf *AuthenticationToken) IsValid() bool {
	return vSelf.Error == nil
}

func (vSelf *AuthenticationToken) Validate() error {

	if vSelf.IsValid() == false {
		return fmt.Errorf("Invalid token %v", vSelf)
	}

	if vSelf.IsExpired() {
		return &TokenExpiredError{}
	}

	return nil
}
