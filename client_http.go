package onedriveclient

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/onedriveclient/auth"
)

var (
	RETRY_MAX      = 5
	RETRY_INTERVAL = time.Second * 1
)

//DoRequest execute an http request to the drive
//parameters:
//  pMethod = http request method (GET,..)
//  pUrl = request url
//  pRequestModifierFunc = optional, func to customize request
//  pResultBean = pointer to a struct where to deserialize data
//returns:
//  error in case of error
func (vSelf *OneDriveClient) DoRequest(pMethod string, pURL string, pRequestModifierFunc func(*http.Request) error, pResultBean interface{}) (vRisError error) {
	var vURL string
	if strings.Contains(pURL, "://") {
		vURL = pURL
	} else {
		vURL = "https://api.onedrive.com/v1.0" + pURL
	}

	if diagnostic.IsLogTrace() {
		diagnostic.LogTrace("OneDriveClient.DoRequest", "Performing request url %s...", vURL)
	}

	vRequest, _ := http.NewRequest(pMethod, vURL, nil)
	if pRequestModifierFunc != nil {
		vModifierError := pRequestModifierFunc(vRequest)
		return vModifierError
	}

	vResponse, vError := vSelf.doRequest(vRequest, 0)
	if vError != nil {
		return diagnostic.NewError("Error performing request", vError)
	}

	vData, vParseDataError := ioutil.ReadAll(vResponse.Body)
	defer vResponse.Body.Close()

	if vParseDataError != nil {
		return diagnostic.NewError("Error parsing response body", vError)
	}

	if diagnostic.IsLogTrace() {
		diagnostic.LogTrace("OneDriveClient.DoRequest", "Response  %s", vData)
	}

	if pResultBean != nil {
		vError = json.Unmarshal(vData, pResultBean)
		if vError != nil {
			return diagnostic.NewError("Error parsing response: %s", vError, vData)
		}
	}
	return nil

}

//DoRequestDownload execute an http request to download an item
//parameters:
//  pMethod = http request method (GET,..)
//  pUrl = request url
//  pWriter = response body writer
//returns:
//  error in case of error
func (vSelf *OneDriveClient) DoRequestDownload(pMethod string, pURL string, pWriter io.Writer) (vRisError error) {

	var vURL string
	if strings.Contains(pURL, "://") {
		vURL = pURL
	} else {
		vURL = "https://api.onedrive.com/v1.0" + pURL
	}

	vRequest, vError := http.NewRequest(pMethod, vURL, nil)
	if vError != nil {
		return diagnostic.NewError("Error creating request %s", vError, pURL)
	}

	vResponse, vError := vSelf.doRequest(vRequest, 0)
	if vError != nil {
		return diagnostic.NewError("Error executing request %s", vError, pURL)
	}

	_, vError = io.Copy(pWriter, vResponse.Body)
	defer vResponse.Body.Close()

	if vError != nil {
		return diagnostic.NewError("Error downloading %s", vError, pURL)

	}
	return nil
}

func (vSelf *OneDriveClient) doRequest(pRequest *http.Request, pRetryNumber int) (*http.Response, error) {

	if diagnostic.IsLogTrace() {
		diagnostic.LogTrace("OneDriveClient.doRequest", "Performing request %s ...", pRequest.URL)
	}
	vTokenError := vSelf.setAuthorization(pRequest)

	if vTokenError != nil {
		return nil, vTokenError
	}

	vResponse, vError := vSelf.httpClient.Do(pRequest)
	if vError != nil {
		return nil, vError
	}
	vResponseCode := -1

	if vResponse != nil {
		vResponseCode = vResponse.StatusCode
	}

	if vResponseCode != 200 {

		if vResponseCode > 500 && pRetryNumber < RETRY_MAX {
			diagnostic.LogWarning("OneDriveClient.doRequest", "Response code %d, retry...", nil, vResponseCode)
			time.Sleep(RETRY_INTERVAL)
			return vSelf.doRequest(pRequest, pRetryNumber+1)
		}
		return nil, diagnostic.NewError("Url %s ResponseCode %d", nil, pRequest.URL, vResponseCode)
	}

	return vResponse, nil
}

func (vSelf *OneDriveClient) setAuthorization(pRequest *http.Request) error {

	var vAuthenticationToken *auth.AuthenticationToken

	if vSelf.authenticationProvider != nil {
		var vAuthenticationTokenError error
		vAuthenticationToken, vAuthenticationTokenError = vSelf.authenticationProvider.GetAuthenticationToken()

		if vAuthenticationTokenError != nil {
			return diagnostic.NewError("Failed to obtain token", vAuthenticationTokenError)
		}
	}

	if vAuthenticationToken == nil {
		return diagnostic.NewError("Missing authorization token", nil)
	}

	vTokenError := vAuthenticationToken.Validate()
	if vTokenError != nil {
		if _, vIsExpired := vTokenError.(*auth.TokenExpiredError); vIsExpired {

			vApplicationInfo, vApplicationInfoError := vSelf.authenticationProvider.GetApplicationInfo()
			if vApplicationInfoError != nil {
				return diagnostic.NewError("Failed to obtain application info", vApplicationInfoError)
			}

			diagnostic.LogInfo("OneDriveClient.setAuthorization", "Token expired, executing refresh...")
			vNewToken, vError := vAuthenticationToken.Refresh(vApplicationInfo)

			if vError != nil {
				return diagnostic.NewError("failed to refresh token", vError)
			}

			var vAuthenticationProviderInterface interface{} = vSelf.authenticationProvider
			vStaticAuthenticationInfo, vIsStaticIsAuthenticationInfo := vAuthenticationProviderInterface.(auth.StaticAuthenticationInfo)

			if vIsStaticIsAuthenticationInfo {
				vStaticAuthenticationInfo.AuthenticationToken = vNewToken
				diagnostic.LogInfo("OneDriveClient.setAuthorization", "Static token updated")
			}
		} else {
			return diagnostic.NewError("Invalid token", vTokenError)
		}
	}

	pRequest.Header.Set("Authorization", "bearer "+vAuthenticationToken.AccessToken)

	return nil
}
