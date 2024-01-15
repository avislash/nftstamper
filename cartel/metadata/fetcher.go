package metadata

import (
	"fmt"

	"github.com/avislash/nftstamper/lib/metadata"
)

type HoundMetadataFetcher struct {
	metadata.Fetcher[HoundMetadata]
}

type MAYCMetadataFetcher struct {
	metadata.Fetcher[MAYCMetadata]
}

type BAYCMetadataFetcher struct {
	metadata.Fetcher[BAYCMetadata]
}

func NewHoundMetadataFetcher(baseURL string) *HoundMetadataFetcher {
	return &HoundMetadataFetcher{metadata.NewJSONMetadataFetcher[HoundMetadata](baseURL)}
}

func NewMAYCMetadataFetcher(baseURL string) *MAYCMetadataFetcher {
	return &MAYCMetadataFetcher{metadata.NewJSONMetadataFetcher[MAYCMetadata](baseURL)}
}

func NewBAYCMetadataFetcher(baseURL string, baseHash string) (*BAYCMetadataFetcher, error) {
	ipfsFetcher, err := metadata.NewIPFSMetadataFetcher[BAYCMetadata](baseURL, baseHash)
	if err != nil {
		return nil, fmt.Errorf("Error creating BAYC Metadata Fetcher: %w", err)
	}
	return &BAYCMetadataFetcher{ipfsFetcher}, nil
}
