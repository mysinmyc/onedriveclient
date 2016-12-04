package onedriveclient

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/mysinmyc/onedriveclient/auth"
)

//OneDriveClient client instance
type OneDriveClient struct {
	authenticationToken *auth.AuthenticationToken
	httpClient          *http.Client
	LastResponse        struct {
		StatusCode int
		Body       string
	}
}

//NewOneDriveClient create a new instance of onedriveclient
func NewOneDriveClient() *OneDriveClient {

	vRis := &OneDriveClient{httpClient: &http.Client{}}

	return vRis
}

//SetAuthenticationToken set current authentication token
func (vSelf *OneDriveClient) SetAuthenticationToken(pAuthenticationToken *auth.AuthenticationToken) *OneDriveClient {
	vSelf.authenticationToken = pAuthenticationToken
	return vSelf
}

//PerformNewAuthentication Execute the interactive flow for a new authentication
func (vSelf *OneDriveClient) PerformNewAuthentication(pAuthenticationHelper auth.AuthenticationHelper, pTimeout time.Duration) error {
	vAuthenticationToken, vError := pAuthenticationHelper.WaitAuthenticationToken(pTimeout)
	vSelf.SetAuthenticationToken(vAuthenticationToken)
	return vError
}

func (vSelf *OneDriveClient) doRequest(pURL string, pResultBean interface{}) (vRisError error) {

	var vURL string
	if strings.Contains(pURL, "://") {
		vURL = pURL
	} else {
		vURL = "https://api.onedrive.com/v1.0" + pURL
	}

	vRequest, _ := http.NewRequest("GET", vURL, nil)
	vRequest.Header.Set("Content-Type", "application/json")

	vTokenError := vSelf.setAuthenticationCode(vRequest)
	if vTokenError != nil {
		return vTokenError
	}

	vResponse, vError := vSelf.httpClient.Do(vRequest)
	vResponseCode := -1
	if vResponse != nil {
		vResponseCode = vResponse.StatusCode
	}
	vSelf.LastResponse.StatusCode = vResponseCode
	vSelf.LastResponse.Body = ""

	if vError != nil || vResponseCode != 200 {
		log.Printf("RC: %d, %v", vResponseCode, vError)
		return vError
	}

	vData, _ := ioutil.ReadAll(vResponse.Body)
	vSelf.LastResponse.Body = string(vData)
	if pResultBean != nil {
		vError = json.Unmarshal(vData, pResultBean)
		if vError != nil {
			log.Printf("ERROR PARSING RESPONSE: %s %v", vData, vError)
			return vError
		}
	}
	return nil
}

func (vSelf *OneDriveClient) doGet(pURL string, pWriter io.Writer) (vRisError error) {

	var vURL string
	if strings.Contains(pURL, "://") {
		vURL = pURL
	} else {
		vURL = "https://api.onedrive.com/v1.0" + pURL
	}

	vRequest, _ := http.NewRequest("GET", vURL, nil)
	vRequest.Header.Set("Content-Type", "data/octetstream")

	vTokenError := vSelf.setAuthenticationCode(vRequest)
	if vTokenError != nil {
		return vTokenError
	}

	vResponse, vError := vSelf.httpClient.Do(vRequest)
	vSelf.LastResponse.StatusCode = vResponse.StatusCode
	vSelf.LastResponse.Body = ""

	if vError != nil || vResponse.StatusCode != 200 {
		log.Printf("RC: %d, %v", vResponse.StatusCode, vError)
		return vError
	}

	_, vError = io.Copy(pWriter, vResponse.Body)

	if vError != nil {
		log.Printf("ERROR DOWNLOADING: %s %v", pURL, vError)
		return vError

	}
	return nil
}

func (vSelf *OneDriveClient) setAuthenticationCode(pRequest *http.Request) error {
	if vSelf.authenticationToken != nil {

		vTokenError := vSelf.authenticationToken.Validate()
		if vTokenError != nil {
			return vTokenError
		}

		pRequest.Header.Set("Authorization", "bearer "+vSelf.authenticationToken.AccessToken)
	}

	return nil
}
