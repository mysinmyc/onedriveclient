package auth

import "time"

//ApplicationInfo returns application informations required for authentication
type ApplicationInfo struct {
	ClientID     string
	ClientSecret string
	Scope        []string
	RedirectURI  string
}

//AuthenticationHelper interface for all the authentication helper
type AuthenticationHelper interface {
	WaitAuthenticationToken(time.Duration) (*AuthenticationToken, error)
	RefreshToken(*AuthenticationToken) (*AuthenticationToken, error)
}
