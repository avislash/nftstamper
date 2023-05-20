package image

import "bytes"

type Combiner interface {
	AdjustImageOpacity(img Image, opacity float64) Image
	CombineImages(img1, img2 Image) Image
	EncodeImage(img Image) (*bytes.Buffer, error)
}
