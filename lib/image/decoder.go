package image

import (
	"image/jpeg"
	"image/png"
	"io"
)

var _ Decoder = (*PNGDecoder)(nil)
var _ Decoder = (*JPEGDecoder)(nil)

type Decoder interface {
	Decode(r io.Reader) (Image, error)
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
