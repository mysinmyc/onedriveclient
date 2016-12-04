package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func reedimCode(pApplicationInfo ApplicationInfo, pCode string) (vRisAuthenticatioToken AuthenticationToken, vRisError error) {
	log.Printf("Asking redeem for authorization code %s, ...", pCode)

	vReedimRequest, _ := http.NewRequest("POST", "https://login.live.com/oauth20_token.srf",
		strings.NewReader(fmt.Sprintf(
			"client_id=%s&redirect_uri=%s&client_secret=%s&code=%s&grant_type=authorization_code",
			pApplicationInfo.ClientID,
			pApplicationInfo.RedirectURI,
			pApplicationInfo.ClientSecret,
			pCode)))

	vReedimRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	vReedimResponse, vReedimError := http.DefaultClient.Do(vReedimRequest)

	if vReedimError != nil {
		log.Printf("ERROR: asking for token reedim:%s", vReedimError)
		vRisError = vReedimError
		return
	}

	if vReedimResponse.StatusCode != http.StatusOK {
		vBodyString, _ := ioutil.ReadAll(vReedimResponse.Body)

		log.Printf("error: %d, %s", vReedimResponse.StatusCode, vBodyString)
		vRisError = fmt.Errorf("error: %d, %s", vReedimResponse.StatusCode, vBodyString)
		return
	}
	vAuthenticationToken := *newAuthenticationToken()

	vDecodeError := json.NewDecoder(vReedimResponse.Body).Decode(&vAuthenticationToken)
	if vDecodeError != nil {
		log.Printf("ERROR: failed to decode token:%s", vDecodeError)
		vRisError = vDecodeError
		return
	}

	return vAuthenticationToken, nil

}
