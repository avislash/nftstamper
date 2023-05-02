package image

import (
	"bytes"
)

type Combiner interface {
	CombineImages(img1, img2 Image) (*bytes.Buffer, error)
}
