package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	_Initialized = false
)

//HttpAuthHelper Helper to authenticate with the drive trough oauth20
// It's supposed to have only one authentication per application
// In start an http server on the specified address
// For more informations about auth workflow check https://dev.onedrive.com/auth/msa_oauth.htm
type HttpAuthHelper struct {
	address               string
	clientID              string
	clientSecret          string
	scope                 []string
	authenticationHandler func(*AuthenticationToken)
}

//NewAuthHelper create a new instance of AuthenticationHelper
func NewHttpAuthHelper(pAddress string, pClientID string, pClientSecret string, pScope []string) *HttpAuthHelper {
	vRis := &HttpAuthHelper{address: pAddress, clientID: pClientID, clientSecret: pClientSecret, scope: pScope}
	vRis.init()
	return vRis
}

//SetAuthenticationHandler Set the function that received AuthenticationTokens coming from authentication flow
func (vSelf *HttpAuthHelper) SetAuthenticationHandler(pAuthenticationHandler func(*AuthenticationToken)) {
	vSelf.authenticationHandler = pAuthenticationHandler
}

func (vSelf *HttpAuthHelper) onAuthenticationToken(pAuthenticationToken *AuthenticationToken) {
	if vSelf.authenticationHandler != nil {
		vSelf.authenticationHandler(pAuthenticationToken)
	}
}

func (vSelf *HttpAuthHelper) onAuthenticationError(pError error) {
	if vSelf.authenticationHandler != nil {
		vSelf.onAuthenticationToken(&AuthenticationToken{Error: pError})
	}
}

func (vSelf *HttpAuthHelper) init() error {

	if _Initialized {
		return nil
	}

	http.HandleFunc("/", func(pResponse http.ResponseWriter, pRequest *http.Request) {

		log.Printf("Asked login, redirecting to microsoft...")
		vMicrosoftLoginURL := fmt.Sprintf(
			"https://login.live.com/oauth20_authorize.srf?client_id=%s&scope=%s&response_type=code&redirect_uri=%s",
			vSelf.clientID,
			strings.Join(vSelf.scope, "%20"),
			url.QueryEscape("http://"+vSelf.address+"/redirect"))
		http.Redirect(pResponse, pRequest, vMicrosoftLoginURL, 302)
	})

	http.HandleFunc("/redirect", func(pResponse http.ResponseWriter, pRequest *http.Request) {

		pRequest.ParseForm()
		vCode := pRequest.FormValue("code")

		log.Printf("Obtained authorization code %s, asking redeem...", vCode)

		vReedimRequest, _ := http.NewRequest("POST", "https://login.live.com/oauth20_token.srf",
			strings.NewReader(fmt.Sprintf(
				"client_id=%s&redirect_uri=%s&client_secret=%s&code=%s&grant_type=authorization_code",
				vSelf.clientID,
				url.QueryEscape("http://"+vSelf.address+"/redirect"),
				vSelf.clientSecret,
				vCode)))

		vReedimRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		vReedimResponse, vReedimError := http.DefaultClient.Do(vReedimRequest)

		if vReedimError != nil {
			log.Printf("ERROR: asking for token reedim:%s", vReedimError)

			vSelf.onAuthenticationError(vReedimError)
			return
		}

		if vReedimResponse.StatusCode != http.StatusOK {
			vBodyString, _ := ioutil.ReadAll(vReedimResponse.Body)

			log.Printf("error: %d, %s", vReedimResponse.StatusCode, vBodyString)
			vSelf.onAuthenticationError(fmt.Errorf("error: %d, %s", vReedimResponse.StatusCode, vBodyString))
			return
		}
		vAuthenticationToken := AuthenticationToken{CreationTime: time.Now()}

		vDecodeError := json.NewDecoder(vReedimResponse.Body).Decode(&vAuthenticationToken)
		if vDecodeError != nil {
			log.Printf("ERROR: failed to decode token:%s", vDecodeError)
			vSelf.onAuthenticationError(vDecodeError)
			return
		}

		log.Printf("reedimResponse token is  %s", vAuthenticationToken.AccessToken)

		vSelf.onAuthenticationToken(&vAuthenticationToken)

		http.Redirect(pResponse, pRequest, "/done", 302)
	})

	http.HandleFunc("/done", func(pResponse http.ResponseWriter, pRequest *http.Request) {
		io.WriteString(pResponse, "<html><body>Authentication succeded</body></html>")
	})

	go http.ListenAndServe(vSelf.address, nil)

	_Initialized = true

	return nil
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
	vSelf.SetAuthenticationHandler(func(pAuthenticationToken *AuthenticationToken) {
		vTokenChannel <- pAuthenticationToken
	})

	select {
	case vToken := <-vTokenChannel:
		return vToken, vToken.Error
	case <-time.After(pTimeout):
		vError := errors.New("timeout expired")
		return &AuthenticationToken{Error: vError}, vError

	}

}
