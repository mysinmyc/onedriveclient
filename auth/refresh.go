package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Refresh refresh a token
//parameters:
//	pApplicationInfo: Informations about applications that request the authentication (must be the same of orignal applciation that obtained the token)
//returns:
//	AuthenticationToken: fresh token
//  error: nil if succeded otherwise the error occurred
func (vSelf *AuthenticationToken) Refresh(pApplicationInfo ApplicationInfo) (vRisAuthenticationToken *AuthenticationToken, vRisError error) {

	log.Println("Asking token refresh...")

	vRefreshRequest, _ := http.NewRequest("POST", "https://login.live.com/oauth20_token.srf",
		strings.NewReader(fmt.Sprintf(
			"client_id=%s&redirect_uri=%s&client_secret=%s&refresh_token=%s&grant_type=refresh_token",
			pApplicationInfo.ClientID,
			pApplicationInfo.RedirectURI,
			pApplicationInfo.ClientSecret,
			vSelf.RefreshToken)))

	vRefreshRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	vRefreshResponse, vRefreshError := http.DefaultClient.Do(vRefreshRequest)

	if vRefreshError != nil {
		log.Printf("ERROR: asking for token refresh:%s", vRefreshError)
		vRisError = vRefreshError
		return
	}

	if vRefreshResponse.StatusCode != http.StatusOK {
		vBodyString, _ := ioutil.ReadAll(vRefreshResponse.Body)

		log.Printf("error: %d, %s", vRefreshResponse.StatusCode, vBodyString)
		vRisError = fmt.Errorf("error: %d, %s", vRefreshResponse.StatusCode, vBodyString)
		return
	}

	vAuthenticationToken := newAuthenticationToken()

	vDecodeError := json.NewDecoder(vRefreshResponse.Body).Decode(vAuthenticationToken)
	if vDecodeError != nil {
		log.Printf("ERROR: failed to decode token:%s\n", vDecodeError)
		vRisError = vDecodeError
		return
	}

	return vAuthenticationToken, nil

}
