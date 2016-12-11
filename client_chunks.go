package onedriveclient

import (
	"net/http"

	"github.com/mysinmyc/gocommons/diagnostic"
)

type Chunk struct {
	Value            []*OneDriveItem `json:"value"`
	ChildrenNextLink string          `json:"@odata.nextLink"`
}

type Chunked interface {
	GetNextChunkLink() string
	GetChildren() []*OneDriveItem
}

func (vSelf *OneDriveItem) GetNextChunkLink() string {
	return vSelf.ChildrenNextLink
}

func (vSelf *OneDriveItem) GetChildren() []*OneDriveItem {
	return vSelf.Children
}

func (vSelf *Chunk) GetNextChunkLink() string {
	return vSelf.ChildrenNextLink
}

func (vSelf *Chunk) GetChildren() []*OneDriveItem {
	return vSelf.Value
}

func (vSelf *OneDriveClient) GetNextChunk(pItem Chunked) (vRisChunk *Chunk, vRisError error) {

	vChunk := &Chunk{}

	if pItem.GetNextChunkLink() == "" {
		return vChunk, nil
	}

	diagnostic.LogDebug("OneDriveClient.GetNextChunk", "Asking next chunk...")
	vError := vSelf.DoRequest(http.MethodGet, pItem.GetNextChunkLink(), nil, vChunk)

	if vError != nil {
		return nil, diagnostic.NewError("Failed to request chunk", vError)
	}

	return vChunk, nil
}

func (vSelf *OneDriveClient) MergeAllChunks(pChunk Chunked) (vRisItems []*OneDriveItem, vRisError error) {

	var vCurChunk Chunked = pChunk

	if vCurChunk.GetNextChunkLink() == "" {

		return vCurChunk.GetChildren(), nil
	}

	diagnostic.LogInfo("OneDriveClient.MergeAllChunks", "response chunked")
	vRis := make([]*OneDriveItem, 0)

	for {

		vNextChunk, vChunkError := vSelf.GetNextChunk(vCurChunk)

		if vChunkError != nil {
			return nil, vChunkError
		}

		for _, vCurSubItem := range vCurChunk.GetChildren() {
			vRis = append(vRis, vCurSubItem)
		}

		if vCurChunk.GetNextChunkLink() == "" {
			break
		}

		vCurChunk = vNextChunk
	}

	return vRis, nil
}
