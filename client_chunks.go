package onedriveclient

import "log"

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

	if pItem.GetNextChunkLink() == "" {
		return &Chunk{}, nil
	}

	vChunk := &Chunk{}

	log.Printf("Asking chunk")
	vError := vSelf.doRequest(pItem.GetNextChunkLink(), vChunk)

	if vError != nil {
		log.Printf("ERROR %v", vError)
		return nil, vError
	}

	return vChunk, nil
}
