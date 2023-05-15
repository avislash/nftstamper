package ipfs

import (
	"strings"

	"github.com/avislash/nftstamper/lib/image"
)

type Client interface {
	GetImageFromIPFS(imagePath string) (image.Image, error)
	setImageDecoder(decoder image.Decoder)
}

type Option = func(c Client)

func WithPNGDecoder() Option {
	return func(c Client) {
		c.setImageDecoder(&image.PNGDecoder{})
	}
}

func WithJPEGDecoder() Option {
	return func(c Client) {
		c.setImageDecoder(&image.JPEGDecoder{})
	}
}

func NewClient(endpoint string, options ...Option) (Client, error) {
	if strings.Contains(endpoint, "https") || strings.Contains(endpoint, "http") {
		return NewWebClient(endpoint, options...)
	}
	return NewIPFSClient(endpoint, options...)
}
