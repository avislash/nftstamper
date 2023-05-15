package ipfs

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/avislash/nftstamper/lib/image"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

var _ Client = (*WebClient)(nil)

type WebClient struct {
	endpoint     string
	ImageDecoder image.Decoder
}

func NewWebClient(endpoint string, options ...Option) (*WebClient, error) {
	_, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error parsing endpoint: %w", err)
	}
	endpoint = strings.TrimSuffix(strings.TrimSuffix(endpoint, "/"), "/ipfs")

	client := &WebClient{endpoint, &image.DefaultDecoder{}}
	for _, applyOpt := range options {
		applyOpt(client)
	}
	return client, nil
}

func (c *WebClient) GetImageFromIPFS(imagePath string) (image.Image, error) {
	// Image CID
	cid := path.New(strings.TrimPrefix(imagePath, "ipfs://"))
	path := c.endpoint + cid.String()

	// Retrieve the file from IPFS
	response, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("Error fetching IPFS Image from %s: %w", path, err)
	}
	defer response.Body.Close()

	img, err := c.ImageDecoder.Decode(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error decoding IPFS File as image: %w", err)
	}

	return img, nil
}

func (c *WebClient) setImageDecoder(decoder image.Decoder) {
	c.ImageDecoder = decoder
}
