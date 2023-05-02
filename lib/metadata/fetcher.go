package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Fetcher[T any] interface {
	Fetch(tokenID uint64) (T, error)
}

type JSONMetadataFetcher[T any] struct {
	baseURL string
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
