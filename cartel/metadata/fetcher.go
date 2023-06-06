package metadata

import (
	"github.com/avislash/nftstamper/lib/metadata"
)

type HoundMetadataFetcher struct {
	metadata.Fetcher[HoundMetadata]
}

type MAYCMetadataFetcher struct {
	metadata.Fetcher[MAYCMetadata]
}

func NewHoundMetadataFetcher(baseURL string) *HoundMetadataFetcher {
	return &HoundMetadataFetcher{metadata.NewJSONMetadataFetcher[HoundMetadata](baseURL)}
}

func NewMAYCMetadataFetcher(baseURL string) *MAYCMetadataFetcher {
	return &MAYCMetadataFetcher{metadata.NewJSONMetadataFetcher[MAYCMetadata](baseURL)}
}
