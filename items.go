package onedriveclient

import "time"

//Item interface for onedrive items
type Item interface {
	GetId() string
	IsFolder() bool
	IsFile() bool
}

//OneDriveItem One drive item object
type OneDriveItem struct {

	//Item Id
	Id string `json:"id"`

	//Item name
	Name string `json:"name"`

	//Web Url
	WebUrl string `json:"webUrl"`

	//Item children (if nil are not loaded)
	Children []*OneDriveItem `json:"children"`

	Folder *struct {
		ChildCount int64 `json:"ChildCount"`
	} `json:"folder"`

	File *struct {
		MimeType           string `json:"mimeType"`
		ProcessingMetadata bool   `json:"processingMetadata"`
		Hashes             struct {
			Crc32Hash    string `json:"Crc32Hash"`
			Sha1Hash     string `json:"Sha1Hash"`
			QuickXorHash string `json:"QuickXorHash"`
		} `json:"hashes"`
	} `json:"file"`

	CreatedBy struct {
		User         *Identity `json:"user"`
		Application  *Identity `json:"application"`
		Devide       *Identity `json:"device"`
		Organization *Identity `json:"organization"`
	} `json:"CreatedBy"`

	CreatedDateTime time.Time `json:"createdDateTime"`

	//If present, contains the url of the next chunk contains the rest of the children
	ChildrenNextLink string `json:"children@odata.nextLink"`
}

type Identity struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

func (vSelf *OneDriveItem) GetId() string {
	return vSelf.Id
}

func (vSelf *OneDriveItem) IsFolder() bool {
	return vSelf.Folder != nil
}

func (vSelf *OneDriveItem) IsFile() bool {
	return vSelf.File != nil
}
