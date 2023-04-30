package ipfs

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"strings"

	"github.com/ipfs/boxo/files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	ipfsClient "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

type Client struct {
	*httpapi.HttpApi
}

func NewClient() (*Client, error) {
	client, err := ipfsClient.NewLocalApi()
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %w", err)
	}
	return &Client{client}, nil
}

func (c *Client) GetSentinelFromIPFS(imagePath string) (image.Image, error) {
	// Sentinel CID
	cid := path.New(strings.TrimPrefix(imagePath, "ipfs://"))

	// Retrieve the file from IPFS
	node, err := c.Unixfs().Get(context.Background(), cid)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving centinel from IPFS Hash %s: %w", cid, err)
	}

	file := files.ToFile((node))
	defer file.Close()

	sentinel, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Error decoding IPFS File as PNG: %w", err)
	}

	return sentinel, nil
}
