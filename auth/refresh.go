package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mysinmyc/gocommons/diagnostic"
)

//Refresh refresh a token
//parameters:
//	pApplicationInfo: Informations about applications that request the authentication (must be the same of orignal applciation that obtained the token)
//returns:
//	AuthenticationToken: fresh token
//  error: nil if succeded otherwise the error occurred
func (vSelf *AuthenticationToken) Refresh(pApplicationInfo ApplicationInfo) (vRisAuthenticationToken *AuthenticationToken, vRisError error) {

	diagnostic.LogDebug("AuthenticationToken.Refresh", "Started token refresh...")

	vRefreshRequest, vRefreshRequestError := http.NewRequest("POST", "https://login.live.com/oauth20_token.srf",
		strings.NewReader(fmt.Sprintf(
			"client_id=%s&redirect_uri=%s&client_secret=%s&refresh_token=%s&grant_type=refresh_token",
			pApplicationInfo.ClientID,
			pApplicationInfo.RedirectURI,
			pApplicationInfo.ClientSecret,
			vSelf.RefreshToken)))

	if vRefreshRequestError != nil {
		return nil, diagnostic.NewError("Failed to create request ", vRefreshRequestError)
	}

	vRefreshRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	vRefreshResponse, vRefreshError := http.DefaultClient.Do(vRefreshRequest)

	if vRefreshError != nil {
		return nil, diagnostic.NewError("Token refresh failed ", vRefreshRequestError)
	}

	if vRefreshResponse.StatusCode != http.StatusOK {
		vBodyString, _ := ioutil.ReadAll(vRefreshResponse.Body)
		return nil, diagnostic.NewError("token request error: %d, %s", nil, vRefreshResponse.StatusCode, vBodyString)
	}

	vDecodeError := json.NewDecoder(vRefreshResponse.Body).Decode(vSelf)
	if vDecodeError != nil {
		return nil, diagnostic.NewError("Failed to decode token ", vDecodeError)
	}
	vSelf.CreationTime = time.Now()

	diagnostic.LogDebug("AuthenticationToken.Refresh", "Token successfully refreshed")
	return vSelf, nil

}
