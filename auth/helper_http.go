package auth

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mysinmyc/gocommons/diagnostic"
)

var (
	_Initialized = false
)

//HttpAuthHelper Helper to authenticate with the drive trough oauth20
// It's supposed to have only one authentication per application
// In start an http server on the specified address
// For more informations about auth workflow check https://dev.onedrive.com/auth/msa_oauth.htm
type HttpAuthHelper struct {
	applicationInfo             ApplicationInfo
	address                     string
	pathPrefix                  string
	authenticationHandler       AuthenticationHandler
	redirectAfterAuthentication string
}

//NewAuthHelper create a new instance of AuthenticationHelper
func NewHttpAuthHelper(pAddress string, pClientID string, pClientSecret string, pScope []string, pRedirectAfterAuthentication string) *HttpAuthHelper {
	vRis := &HttpAuthHelper{address: pAddress, applicationInfo: ApplicationInfo{ClientID: pClientID, ClientSecret: pClientSecret, Scope: pScope, RedirectURI: "http://" + pAddress + "/onedrive/auth/redirect"}, redirectAfterAuthentication: pRedirectAfterAuthentication}
	vRis.init()
	return vRis
}

//SetAuthenticationHandler Set the function that received AuthenticationTokens coming from authentication flow
func (vSelf *HttpAuthHelper) SetAuthenticationHandler(pAuthenticationHandler AuthenticationHandler) {
	vSelf.authenticationHandler = pAuthenticationHandler
}

func (vSelf *HttpAuthHelper) onAuthenticationToken(pAuthenticationToken *AuthenticationToken) {
	if vSelf.authenticationHandler != nil {
		vSelf.authenticationHandler(pAuthenticationToken, vSelf.applicationInfo)
	}
}

func (vSelf *HttpAuthHelper) onAuthenticationError(pError error) {
	if vSelf.authenticationHandler != nil {
		vSelf.onAuthenticationToken(newAuthenticationTokenError(pError))
	}
}

func (vSelf *HttpAuthHelper) init() error {

	if _Initialized {
		return nil
	}

	http.HandleFunc("/onedrive/auth/begin", func(pResponse http.ResponseWriter, pRequest *http.Request) {

		io.WriteString(pResponse, fmt.Sprintf(`

		<html>
		<body>

		<br/><br/>
		
		<center>


		<form method="POST" action="redirectToMicrosoft" target="_top">
		<table>
		<tr><td>Client id</td><td><input name="client_id" type="text" value="%s"/></td></tr>
		<tr><td>Client secret</td><td><input name="client_secret" type="password" value="%s"/></td></tr>
		<tr><td>&nbsp;</td><td><input type="submit" value="Begin authentication"/></td></tr>
		</table>		
		</form>
		</center>
		
		<br/><br/><br/><br/>
		To perform authentication you must register your application on <a href="https://apps.dev.microsoft.com/" target="_new"/>Microsoft Application Registration Portal</a>.<br/>Application redirect is %s 


		
		</body>
		</html>
		`, vSelf.applicationInfo.ClientID, vSelf.applicationInfo.ClientSecret, vSelf.applicationInfo.RedirectURI))
	})

	http.HandleFunc("/onedrive/auth/redirectToMicrosoft", func(pResponse http.ResponseWriter, pRequest *http.Request) {

		vClientId := pRequest.FormValue("client_id")
		if vClientId != "" {
			vSelf.applicationInfo.ClientID = vClientId
		}

		vClientSecret := pRequest.FormValue("client_secret")
		if vClientSecret != "" {
			vSelf.applicationInfo.ClientSecret = vClientSecret
		}

		diagnostic.LogInfo("HttpAuthHelper", "Asked login, redirecting to microsoft...")
		vMicrosoftLoginURL := fmt.Sprintf(
			"https://login.live.com/oauth20_authorize.srf?client_id=%s&scope=%s&response_type=code&redirect_uri=%s",
			vSelf.applicationInfo.ClientID,
			strings.Join(vSelf.applicationInfo.Scope, "%20"),
			url.QueryEscape(vSelf.applicationInfo.RedirectURI))
		http.Redirect(pResponse, pRequest, vMicrosoftLoginURL, 302)
	})

	http.HandleFunc("/onedrive/auth/redirect", func(pResponse http.ResponseWriter, pRequest *http.Request) {

		pRequest.ParseForm()
		vCode := pRequest.FormValue("code")

		diagnostic.LogInfo("Asking for token reedim authorization code %s, asking redeem...", vCode)

		vAuthenticationToken, vReedimError := reedimCode(vSelf.applicationInfo, vCode)

		if vReedimError != nil {
			vSelf.onAuthenticationError(vReedimError)
			return
		}

		vSelf.onAuthenticationToken(&vAuthenticationToken)

		http.Redirect(pResponse, pRequest, "/onedrive/auth/done", 302)
	})

	http.HandleFunc("/onedrive/auth/done", func(pResponse http.ResponseWriter, pRequest *http.Request) {

		if vSelf.redirectAfterAuthentication != "" {
			http.Redirect(pResponse, pRequest, vSelf.redirectAfterAuthentication, 302)
			return
		}
		io.WriteString(pResponse, `<html><body>


			<script>
				setTimeout(100)
			</script>
			
			Authentication succeded

			
			</body></html>
			`)
	})

	_Initialized = true

	return nil
}

//StartListener method required in case authentication helper is not bound to an existing httpd
func (vSelf *HttpAuthHelper) StartListener() {
	go http.ListenAndServe(vSelf.address, nil)
}

//WaitAuthenticationToken block the current thread waiting for a token coming from authentication flow
// parameters:
//		pTimeout timeout interval
// returns:
//		vRisAuthenticationToken authenticationToken generated. In case of errors, is not valid and contains in the field AuthenticationToken.Error the cause
//		vRisError nil in case of authentication succeded otherwise the error occurred
func (vSelf *HttpAuthHelper) WaitAuthenticationToken(pTimeout time.Duration) (vRisAuthenticationToken *AuthenticationToken, vRisError error) {

	log.Printf("Requested authentication, open a browser to the following url to continue: http://%s", vSelf.address)
	vTokenChannel := make(chan *AuthenticationToken)
	vSelf.SetAuthenticationHandler(func(pAuthenticationToken *AuthenticationToken, pApplicationInfo ApplicationInfo) {
		vTokenChannel <- pAuthenticationToken
	})

	select {
	case vToken := <-vTokenChannel:
		return vToken, vToken.Error
	case <-time.After(pTimeout):
		vError := diagnostic.NewError("timeout expired waiting for authentication", nil)
		return newAuthenticationTokenError(vError), vError

	}

}

func (vSelf *HttpAuthHelper) RefreshToken(pAuthenticationToken *AuthenticationToken) (vRisToken *AuthenticationToken, vRisError error) {
	vRisToken, vRisError = pAuthenticationToken.Refresh(vSelf.applicationInfo)
	return
}

func (vSelf *HttpAuthHelper) GetApplicationInfo() ApplicationInfo {
	return vSelf.applicationInfo
}
