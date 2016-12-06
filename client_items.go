package onedriveclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

//GetChildren get item children
func (vSelf *OneDriveClient) GetChildren(pItem interface{}) (vRisChildren *Chunk, vRisError error) {

	vRis := &Chunk{}

	switch pItem.(type) {

	case string:
		vRisError = vSelf.DoRequest(http.MethodGet, "/drive/root:/"+pItem.(string)+"/children", nil, vRis)

	case *Drive:
		vRisError = vSelf.DoRequest(http.MethodGet, "/drives/"+pItem.(*Drive).GetId()+"/root/children", nil, vRis)

	case *OneDriveItem:
		vRisError = vSelf.DoRequest(http.MethodGet, "/drive/items/"+pItem.(Item).GetId()+"/children", nil, vRis)

	default:
		vRisError = fmt.Errorf("Invalid item type %T", pItem)
	}

	if vRisError != nil {
		log.Printf("Error getting children of %v: %v", pItem, vRisError)
	}

	return vRis, nil
}

//GetItem Item metadata
func (vSelf *OneDriveClient) GetItem(pItem interface{}, pExpandChildren bool) (vRisItem *OneDriveItem, vRisError error) {

	vItemData := &OneDriveItem{}

	vSuffix := ""
	if pExpandChildren {
		vSuffix = "?expand=children"
	}

	switch pItem.(type) {

	case string:
		vRisError = vSelf.DoRequest(http.MethodGet, "/drive/root:/"+pItem.(string)+vSuffix, nil, vItemData)
	case *Drive:
		vRisError = vSelf.DoRequest(http.MethodGet, "/drives/"+pItem.(*Drive).GetId()+"/root"+vSuffix, nil, vItemData)
	case *OneDriveItem:
		vRisError = vSelf.DoRequest(http.MethodGet, "/drive/items/"+pItem.(Item).GetId()+vSuffix, nil, vItemData)
	default:
		vRisError = fmt.Errorf("Invalid item type %T", pItem)

	}

	if vRisError != nil {
		log.Printf("Error getting item of %v: %v", pItem, vRisError)
	}

	return vItemData, nil
}

//DownloadContentInto download an item content
func (vSelf *OneDriveClient) DownloadContentInto(pItem *OneDriveItem, pWriter io.Writer) error {
	return vSelf.DoRequestDownload(http.MethodGet, "/drive/items/"+pItem.GetId()+"/content", pWriter)
}
