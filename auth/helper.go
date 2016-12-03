package auth

import (
	"time"
)

//AuthenticationHelper interface for all the authentication helper
type AuthenticationHelper interface {
	WaitAuthenticationToken(time.Duration) (*AuthenticationToken, error)
}
