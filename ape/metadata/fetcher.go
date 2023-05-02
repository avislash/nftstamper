package metadata

import (
	"github.com/avislash/nftstamper/lib/metadata"
)

type SentinelMetadataFetcher struct {
	metadata.Fetcher[SentinelMetadata]
}

func NewSentinelMetadataFetcher(baseURL string) *SentinelMetadataFetcher {
	return &SentinelMetadataFetcher{metadata.NewJSONMetadataFetcher[SentinelMetadata](baseURL)}
}
