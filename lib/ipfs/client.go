package ipfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/avislash/nftstamper/lib/image"
	"github.com/ipfs/boxo/files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	ipfsClient "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

type Client struct {
	*httpapi.HttpApi
	ImageDecoder image.Decoder
}

func NewClient(imageDecoder image.Decoder) (*Client, error) {
	client, err := ipfsClient.NewLocalApi()
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %w", err)
	}
	return &Client{client, imageDecoder}, nil
}

func (c *Client) GetImageFromIPFS(imagePath string) (image.Image, error) {
	// Image CID
	cid := path.New(strings.TrimPrefix(imagePath, "ipfs://"))

	// Retrieve the file from IPFS
	node, err := c.Unixfs().Get(context.Background(), cid)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving centinel from IPFS Hash %s: %w", cid, err)
	}

	file := files.ToFile((node))
	defer file.Close()

	img, err := c.ImageDecoder.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Error decoding IPFS File as PNG: %w", err)
	}

	return img, nil
}
