package onedriveclient

import "log"
import "fmt"

//GetChildren get item children
func (vSelf *OneDriveClient) GetChildren(pItem interface{}) (vRisChildren *Chunk, vRisError error) {

	vRis := &Chunk{}

	switch pItem.(type) {
	case *Drive:
		vRisError = vSelf.doRequest("/drives/"+pItem.(*Drive).GetId()+"/root/children", vRis)

	case *OneDriveItem:
		vRisError = vSelf.doRequest("/drive/items/"+pItem.(Item).GetId()+"/children", vRis)

	default:
		vRisError = fmt.Errorf("Invalid item type %T", pItem)
	}

	if vRisError != nil {
		log.Printf("Error getting children of %v: %v", pItem, vRisError)
	}

	return vRis, nil
}

//GetItemByPath get an item by his path
func (vSelf *OneDriveClient) GetItemByPath(pPath string, pExpandChildren bool) (vRisItem *OneDriveItem, vRisError error) {

	vItemData := &OneDriveItem{}

	vSuffix := ""
	if pExpandChildren {
		vSuffix = "?expand=children"
	}
	vRisError = vSelf.doRequest("/drive/root:/"+pPath+vSuffix, vItemData)

	if vRisError != nil {
		log.Printf("Error getting item by path %v: %v", pPath, vRisError)
	}

	return vItemData, nil
}

//GetItem Item metadata
func (vSelf *OneDriveClient) GetItem(pItem interface{}, pExpandChildren bool) (vRisItem *OneDriveItem, vRisError error) {

	vItemData := &OneDriveItem{}

	vSuffix := ""
	if pExpandChildren {
		vSuffix = "?expand=children"
	}

	switch pItem.(type) {
	case *Drive:
		vRisError = vSelf.doRequest("/drives/"+pItem.(*Drive).GetId()+"/root"+vSuffix, vItemData)
	case *OneDriveItem:
		vRisError = vSelf.doRequest("/drive/items/"+pItem.(Item).GetId()+vSuffix, vItemData)
	default:
		vRisError = fmt.Errorf("Invalid item type %T", pItem)

	}

	if vRisError != nil {
		log.Printf("Error getting item of %v: %v", pItem, vRisError)
	}

	return vItemData, nil
}
