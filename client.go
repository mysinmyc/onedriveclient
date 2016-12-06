package onedriveclient

import (
	"net/http"

	"time"

	"github.com/mysinmyc/onedriveclient/auth"
)

//OneDriveClient client instance
type OneDriveClient struct {
	authenticationToken *auth.AuthenticationToken
	applicationInfo     auth.ApplicationInfo
	httpClient          *http.Client
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

//SetApplicationInfo set application info to refresh token
func (vSelf *OneDriveClient) SetApplicationInfo(pApplicationInfo auth.ApplicationInfo) *OneDriveClient {
	vSelf.applicationInfo = pApplicationInfo
	return vSelf
}

//PerformNewAuthentication Execute the interactive flow for a new authentication
func (vSelf *OneDriveClient) PerformNewAuthentication(pAuthenticationHelper auth.AuthenticationHelper, pTimeout time.Duration) error {
	vAuthenticationToken, vError := pAuthenticationHelper.WaitAuthenticationToken(pTimeout)
	if vError != nil {
		return vError
	}
	vSelf.SetAuthenticationToken(vAuthenticationToken)
	vSelf.SetApplicationInfo(pAuthenticationHelper.GetApplicationInfo())
	return nil
}
