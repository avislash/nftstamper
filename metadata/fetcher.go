package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SentinelMetadataFetcher struct {
	baseURL string
}

func NewSentinelMetadataFetcher(baseURL string) *SentinelMetadataFetcher {
	return &SentinelMetadataFetcher{baseURL}
}

func (smf *SentinelMetadataFetcher) FetchMetdata(tokenID uint64) (SentinelMetadata, error) {
	var metadata SentinelMetadata
	response, err := http.Get(fmt.Sprintf("%s/%d", smf.baseURL, tokenID))
	if err != nil {
		return SentinelMetadata{}, fmt.Errorf("Error fetching metadata for token %d: %w", tokenID, err)
	}

	err = json.NewDecoder(response.Body).Decode(&metadata)
	if err != nil {
		return SentinelMetadata{}, fmt.Errorf("Error unmarshalling metadata for token %d: %w", tokenID, err)
	}
	return metadata, nil
}
