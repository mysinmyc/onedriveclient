package onedriveclient

import (
	"net/http"

	"github.com/mysinmyc/onedriveclient/auth"
)

//OneDriveClient client instance
type OneDriveClient struct {
	authenticationProvider auth.AuthenticationProvider
	httpClient             *http.Client
}

//NewOneDriveClient create a new instance of onedriveclient
func NewOneDriveClient() *OneDriveClient {

	vRis := &OneDriveClient{httpClient: &http.Client{}}

	return vRis
}
