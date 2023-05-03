package image

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
)

var _ Decoder = (*PNGDecoder)(nil)
var _ Decoder = (*JPEGDecoder)(nil)

type Decoder interface {
	Decode(r io.Reader) (Image, error)
}

type DefaultDecoder struct {
}

func (dd *DefaultDecoder) Decode(r io.Reader) (Image, error) {
	return nil, fmt.Errorf("Default Decoder unable to decode")
}

type PNGDecoder struct {
}

func (pd *PNGDecoder) Decode(r io.Reader) (Image, error) {
	return png.Decode(r)
}

type JPEGDecoder struct {
}

func (jd *JPEGDecoder) Decode(r io.Reader) (Image, error) {
	return jpeg.Decode(r)
}
