package onedriveclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mysinmyc/gocommons/diagnostic"
)

//GetChildren get item children
func (vSelf *OneDriveClient) GetChildren(pItem interface{}) (vRisChildren *Chunk, vRisError error) {

	vRis := &Chunk{}

	vItemPath, vPathError := vSelf.GetItemPath(pItem)

	if vPathError != nil {
		return nil, vPathError
	}
	vRisError = vSelf.DoRequest(http.MethodGet, vItemPath+"/children", nil, vRis)

	if vRisError != nil {
		log.Printf("Error getting children of %v: %v", pItem, vRisError)
	}

	return vRis, nil
}

func escapeChar(pChar byte) []byte {

	switch pChar {
	case ' ':
		return []byte("%20")
	default:
		return []byte{pChar}
	}
}
func encodePath(pPath string) string {

	return pPath

	/*
		vRis := make([]byte, 0, 10)

		for vCnt := 0; vCnt < len(pPath); vCnt++ {

			vRis = append(vRis, escapeChar(pPath[vCnt])...)
		}

		return string(vRis)
	*/
}

func (vSelf *OneDriveClient) GetItemPath(pItem interface{}) (vRisPath string, vRisError error) {

	switch pItem.(type) {
	case string:
		vRisPath = pItem.(string)

		if strings.Contains(vRisPath, "/") == false {
			vRisPath = "/drive/items/" + vRisPath
		} else {
			if strings.HasPrefix(vRisPath, "/drive") == false {
				vRisPath = "/drive/root:/" + encodePath(vRisPath)
			}
		}
	case *Drive:
		vRisPath = "/drives/" + pItem.(*Drive).GetId() + "/root"
	case *OneDriveItem:
		vRisPath = "/drive/items/" + pItem.(Item).GetId()
	default:
		vRisError = fmt.Errorf("Invalid item type %T", pItem)

	}

	if vRisPath != "" {
		vRisPath = strings.Replace(vRisPath, "//", "/", -1)
	}
	return
}

//GetItem Item metadata
func (vSelf *OneDriveClient) GetItem(pItem interface{}, pExpandChildren bool) (vRisItem *OneDriveItem, vRisError error) {

	vItemData := &OneDriveItem{}

	vSuffix := ""
	if pExpandChildren {
		vSuffix = "?expand=children"
	}

	vItemPath, vPathError := vSelf.GetItemPath(pItem)

	if vPathError != nil {
		return nil, vPathError
	}
	vRisError = vSelf.DoRequest(http.MethodGet, vItemPath+vSuffix, nil, vItemData)

	if vRisError != nil {
		return nil, diagnostic.NewError("Failed", vRisError)
	}

	if vItemData.GetNextChunkLink() != "" {

		vChildren, vChunkError := vSelf.MergeAllChunks(vItemData)

		if vChunkError != nil {
			return nil, vChunkError
		}

		vItemData.Children = vChildren

		if vItemData.Folder != nil {
			if int64(len(vItemData.Children)) != vItemData.Folder.ChildCount {

				return nil, fmt.Errorf("%d children expected for item %v, effective %d ", vItemData.Folder.ChildCount, vItemData.Id, len(vItemData.Children))

			}
		}
	}

	return vItemData, nil
}

//DownloadContentInto download an item content
func (vSelf *OneDriveClient) DownloadContentInto(pItem interface{}, pWriter io.Writer) error {

	vItemPath, vPathError := vSelf.GetItemPath(pItem)

	if vPathError != nil {
		return vPathError
	}

	return vSelf.DoRequestDownload(http.MethodGet, vItemPath+"/content", pWriter)
}
