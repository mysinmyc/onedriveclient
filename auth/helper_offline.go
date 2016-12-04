package auth

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

//OfflineAuthHelper Helper to authenticate with the drive trough oauth20 without an http server
// Application must be registered for mobile access!!!
// It requires a manual user setup
// For more informations about auth workflow check https://dev.onedrive.com/auth/msa_oauth.htm
type OfflineAuthHelper struct {
	applicationInfo       ApplicationInfo
	authenticationHandler func(*AuthenticationToken)
}

//NewOfflineAuthHelper create a new instance of OfflineAuthHelper
func NewOfflineAuthHelper(pClientID string, pClientSecret string, pScope []string) *OfflineAuthHelper {
	vRis := &OfflineAuthHelper{applicationInfo: ApplicationInfo{ClientID: pClientID, ClientSecret: pClientSecret, Scope: pScope, RedirectURI: "https://login.live.com/oauth20_desktop.srf"}}
	return vRis
}

//SetAuthenticationHandler Set the function that received AuthenticationTokens coming from authentication flow
func (vSelf *OfflineAuthHelper) SetAuthenticationHandler(pAuthenticationHandler func(*AuthenticationToken)) {
	vSelf.authenticationHandler = pAuthenticationHandler
}

func (vSelf *OfflineAuthHelper) onAuthenticationToken(pAuthenticationToken *AuthenticationToken) {
	if vSelf.authenticationHandler != nil {
		vSelf.authenticationHandler(pAuthenticationToken)
	}
}

func (vSelf *OfflineAuthHelper) onAuthenticationError(pError error) {
	if vSelf.authenticationHandler != nil {
		vSelf.onAuthenticationToken(&AuthenticationToken{Error: pError})
	}
}

func (vSelf *OfflineAuthHelper) asyncReadTokenFromStdin() {
	log.Printf("Open the browser to the following url %s", vSelf.GetAuthorizationURL())
	vRedirectURI := ""
	log.Print("Then paste the result url here: ")
	fmt.Scanf("%s\n", &vRedirectURI)
	vReedimToken, _ := vSelf.ReedimTokenFromRedirectURI(vRedirectURI)
	vSelf.onAuthenticationToken(vReedimToken)
}

//WaitAuthenticationToken block the current thread waiting for a token coming from authentication flow
// parameters:
//		pTimeout timeout interval
// returns:
//		vRisAuthenticationToken authenticationToken generated. In case of errors, is not valid and contains in the field AuthenticationToken.Error the cause
//		vRisError nil in case of authentication succeded otherwise the error occurred
func (vSelf *OfflineAuthHelper) WaitAuthenticationToken(pTimeout time.Duration) (vRisAuthenticationToken *AuthenticationToken, vRisError error) {

	vTokenChannel := make(chan *AuthenticationToken)

	vSelf.SetAuthenticationHandler(func(pAuthenticationToken *AuthenticationToken) {
		vTokenChannel <- pAuthenticationToken
	})

	go vSelf.asyncReadTokenFromStdin()

	select {
	case vToken := <-vTokenChannel:
		return vToken, vToken.Error
	case <-time.After(pTimeout):
		vError := errors.New("timeout expired waiting for authentication")
		return newAuthenticationTokenError(vError), vError

	}

}

//GetAuthorizationURL return the url to be invoked to obtain the authentication token
func (vSelf *OfflineAuthHelper) GetAuthorizationURL() string {
	vMicrosoftLoginURL := fmt.Sprintf(
		"https://login.live.com/oauth20_authorize.srf?client_id=%s&scope=%s&response_type=code&redirect_uri=%s",
		vSelf.applicationInfo.ClientID,
		strings.Join(vSelf.applicationInfo.Scope, "%20"),
		url.QueryEscape(vSelf.applicationInfo.RedirectURI))

	return vMicrosoftLoginURL
}

//ReedimTokenFromRedirectURI tries to reedim a token from the given redirect_uri
func (vSelf *OfflineAuthHelper) ReedimTokenFromRedirectURI(pURI string) (vRisAuthenticationToken *AuthenticationToken, vRisError error) {

	vURI, vError := url.ParseRequestURI(pURI)

	if vError != nil {
		return newAuthenticationTokenError(vError), vError
	}

	vCode := vURI.Query().Get("code")

	if vCode == "" {
		vError := fmt.Errorf("missing code parameter in uri %s ", pURI)
		return newAuthenticationTokenError(vError), vError
	}

	vAuthenticationToken, vReedimError := reedimCode(vSelf.applicationInfo, vCode)

	if vReedimError != nil {
		vSelf.onAuthenticationError(vReedimError)
		return newAuthenticationTokenError(vReedimError), vReedimError
	}

	vSelf.onAuthenticationToken(&vAuthenticationToken)
	return &vAuthenticationToken, nil
}

func (vSelf *OfflineAuthHelper) RefreshToken(pAuthenticationToken *AuthenticationToken) (vRisToken *AuthenticationToken, vRisError error) {
	vRisToken, vRisError = pAuthenticationToken.Refresh(vSelf.applicationInfo)
	return
}
