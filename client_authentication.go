package onedriveclient

import (
	"time"

	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/onedriveclient/auth"
)

//SetAuthenticationInfo set current authentication informations
func (vSelf *OneDriveClient) SetAuthenticationProvider(pAuthenticationProvider auth.AuthenticationProvider) *OneDriveClient {
	vSelf.authenticationProvider = pAuthenticationProvider
	return vSelf
}

//PerformNewAuthentication Execute the interactive flow for a new authentication
func (vSelf *OneDriveClient) PerformNewAuthentication(pAuthenticationHelper auth.AuthenticationHelper, pTimeout time.Duration) error {
	vAuthenticationToken, vError := pAuthenticationHelper.WaitAuthenticationToken(pTimeout)
	if vError != nil {
		return diagnostic.NewError("Error performing new authentication", vError)
	}
	vApplicationInfo := pAuthenticationHelper.GetApplicationInfo()
	vSelf.SetAuthenticationProvider(&auth.StaticAuthenticationInfo{AuthenticationToken: vAuthenticationToken, ApplicationInfo: &vApplicationInfo})
	return nil
}
