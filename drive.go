package onedriveclient

import (
	"net/http"
)

//Quota contains drive's quota statics
type Quota struct {
	Total     int    `json:"total"`
	Used      int    `json:"used"`
	Remaining int    `json:"remaining"`
	Deleted   int    `json:"deleted"`
	State     string `json:"state"`
}

//Drive
type Drive struct {
	Id        string `json:"id"`
	DriveType string `json:"driveType"`
	Quota     Quota  `json:"quota"`
}

//GetDefaultDrive return the default Drive
func (vSelf *OneDriveClient) GetDefaultDrive() (vRisDrive *Drive, vRisError error) {

	vRis := &Drive{}
	vRisError = vSelf.DoRequest(http.MethodGet, "/drive", nil, vRis)
	vRisDrive = vRis
	return
}

func (vSelf *Drive) GetId() string {
	return vSelf.Id
}
