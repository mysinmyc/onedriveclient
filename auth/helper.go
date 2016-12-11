package auth

import "time"

type AuthenticationHandler func(*AuthenticationToken, ApplicationInfo)

//AuthenticationHelper interface for all the authentication helper
type AuthenticationHelper interface {
	WaitAuthenticationToken(time.Duration) (*AuthenticationToken, error)
	RefreshToken(*AuthenticationToken) (*AuthenticationToken, error)
	GetApplicationInfo() ApplicationInfo
	SetAuthenticationHandler(pAuthenticationHandler AuthenticationHandler)
}
