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
	multiAddr "github.com/multiformats/go-multiaddr"
)

type Client struct {
	*httpapi.HttpApi
	ImageDecoder image.Decoder
}

type Option = func(c *Client)

func WithPNGDecoder() Option {
	return func(c *Client) {
		c.ImageDecoder = &image.PNGDecoder{}
	}
}

func WithJPEGDecoder() Option {
	return func(c *Client) {
		c.ImageDecoder = &image.JPEGDecoder{}
	}
}

//endpoint must be in MultiAddr Format as specified under https://github.com/multiformats/multiaddr#encoding
func NewClient(endpoint string, options ...Option) (*Client, error) {
	addr, err := multiAddr.NewMultiaddr(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %w", err)
	}

	client, err := ipfsClient.NewApi(addr)
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %w", err)
	}

	c := &Client{client, &image.DefaultDecoder{}}

	for _, applyOpt := range options {
		applyOpt(c)
	}

	return c, nil
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
		return nil, fmt.Errorf("Error decoding IPFS File as image: %w", err)
	}

	return img, nil
}
