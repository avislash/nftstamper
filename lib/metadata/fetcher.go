package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/avislash/nftstamper/lib/ipfs"
)

type Fetcher[T any] interface {
	Fetch(tokenID uint64) (T, error)
}

type JSONMetadataFetcher[T any] struct {
	baseURL string
}

type IPFSMetadataFetcher[T any] struct {
	ipfsClient ipfs.Client
	baseHash   string
}

func NewJSONMetadataFetcher[T any](baseURL string) *JSONMetadataFetcher[T] {
	return &JSONMetadataFetcher[T]{strings.TrimSuffix(baseURL, "/")}
}

func (jmf *JSONMetadataFetcher[T]) Fetch(tokenID uint64) (T, error) {
	var metadata T
	response, err := http.Get(fmt.Sprintf("%s/%d", jmf.baseURL, tokenID))
	if err != nil {
		return metadata, fmt.Errorf("Error fetching metadata for token %d: %w", tokenID, err)
	}

	err = json.NewDecoder(response.Body).Decode(&metadata)
	if err != nil {
		return metadata, fmt.Errorf("Error unmarshalling metadata for token %d: %w", tokenID, err)
	}
	return metadata, nil
}

func NewIPFSMetadataFetcher[T any](ipfsNodeEndpoint string, baseHash string) (*IPFSMetadataFetcher[T], error) {
	ipfsClient, err := ipfs.NewClient(ipfsNodeEndpoint)
	if err != nil {
		return nil, fmt.Errorf("Error Creating IPFS Metadata Fetcher: %w", err)
	}

	//Ensure bashHash ends with a single /
	baseHash = strings.TrimSuffix(baseHash, "/") + "/"

	return &IPFSMetadataFetcher[T]{ipfsClient, baseHash}, nil
}

func (imf *IPFSMetadataFetcher[T]) Fetch(tokenID uint64) (T, error) {
	var metadata T

	file, err := imf.ipfsClient.GetFileFromIPFS(imf.baseHash + strconv.FormatUint(tokenID, 10))
	if err != nil {
		return metadata, fmt.Errorf("Error fetching Metadata file from IPFS: %w", err)
	}

	err = json.NewDecoder(file).Decode(&metadata)
	if err != nil {
		return metadata, fmt.Errorf("Error unmarshaling Metadata from file: %w", err)
	}

	return metadata, nil
}
