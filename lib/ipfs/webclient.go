package ipfs

import (
	"fmt"
	"io"
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

	file, err := c.GetFileFromIPFS(imagePath)
	defer file.Close()
	img, err := c.ImageDecoder.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Error decoding IPFS File as image: %w", err)
	}

	return img, nil
}

func (c *WebClient) GetFileFromIPFS(filePath string) (io.ReadCloser, error) {
	// CID
	cid := path.New(strings.TrimPrefix(filePath, "ipfs://"))
	path := c.endpoint + cid.String()

	// Retrieve the file from IPFS
	response, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("Error fetching IPFS Image from %s: %w", path, err)
	}

	return response.Body, nil
}

func (c *WebClient) setImageDecoder(decoder image.Decoder) {
	c.ImageDecoder = decoder
}
