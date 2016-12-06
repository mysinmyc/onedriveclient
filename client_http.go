package onedriveclient

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mysinmyc/onedriveclient/auth"
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

	vRequest, _ := http.NewRequest(pMethod, vURL, nil)
	if pRequestModifierFunc != nil {
		vModifierError := pRequestModifierFunc(vRequest)
		return vModifierError
	}

	vResponse, vError := vSelf.doRequest(vRequest)
	if vError != nil {
		return vError
	}

	vData, _ := ioutil.ReadAll(vResponse.Body)
	defer vResponse.Body.Close()

	if pResultBean != nil {
		vError = json.Unmarshal(vData, pResultBean)
		if vError != nil {
			log.Printf("ERROR PARSING RESPONSE: %s %v", vData, vError)
			return vError
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

	vRequest, _ := http.NewRequest(pMethod, vURL, nil)

	vResponse, vError := vSelf.doRequest(vRequest)

	if vError != nil {
		return vError
	}

	_, vError = io.Copy(pWriter, vResponse.Body)
	defer vResponse.Body.Close()

	if vError != nil {
		log.Printf("ERROR DOWNLOADING: %s %v", pURL, vError)
		return vError

	}
	return nil
}

func (vSelf *OneDriveClient) doRequest(pRequest *http.Request) (*http.Response, error) {

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
		return nil, vError
	}

	return vResponse, nil
}

func (vSelf *OneDriveClient) setAuthorization(pRequest *http.Request) error {

	if vSelf.authenticationToken == nil {
		return errors.New("Missing authorization token")
	}

	vTokenError := vSelf.authenticationToken.Validate()
	if vTokenError != nil {
		if _, vIsExpired := vTokenError.(*auth.TokenExpiredError); vIsExpired {
			vNewToken, vError := vSelf.authenticationToken.Refresh(vSelf.applicationInfo)

			if vError != nil {
				return vError
			}

			vSelf.authenticationToken = vNewToken
		} else {
			return vTokenError
		}
	}

	pRequest.Header.Set("Authorization", "bearer "+vSelf.authenticationToken.AccessToken)

	return nil
}
