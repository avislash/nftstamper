package image

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
)

var _ Combiner = (*PNGCombiner)(nil)

type PNGCombinerOption func(p *PNGCombiner)

type PNGCombiner struct {
	pngEncoder png.Encoder
}

func WithCompressionLevel(level png.CompressionLevel) PNGCombinerOption {
	return func(p *PNGCombiner) {
		p.pngEncoder.CompressionLevel = level
	}
}

func NewPNGCombiner(opts ...PNGCombinerOption) *PNGCombiner {
	p := &PNGCombiner{}

	for _, applyOpt := range opts {
		applyOpt(p)
	}

	return p
}

func (pc *PNGCombiner) CombineImages(img1, img2 Image) (*bytes.Buffer, error) {
	// Create a new image with the size of the larger image
	combinedWidth := max(img1.Bounds().Max.X, img2.Bounds().Max.X)
	combinedHeight := max(img1.Bounds().Max.Y, img2.Bounds().Max.Y)
	combinedImg := image.NewRGBA(image.Rect(0, 0, combinedWidth, combinedHeight))

	// Draw the first image onto the combined image
	draw.Draw(combinedImg, img1.Bounds(), img1, image.ZP, draw.Src)

	// Draw the second image onto the combined image with an offset
	offset := image.Pt((combinedWidth-img2.Bounds().Dx())/2, (combinedHeight-img2.Bounds().Dy())/2)
	drawRect := img2.Bounds()
	drawRect = drawRect.Add(offset)
	drawRect = drawRect.Intersect(combinedImg.Bounds())
	drawRect = drawRect.Sub(offset)
	drawRect = drawRect.Add(offset)
	drawRect = drawRect.Intersect(img2.Bounds())
	drawRect = drawRect.Sub(offset)
	draw.Draw(combinedImg, drawRect, img2, img2.Bounds().Min, draw.Over)

	buff := new(bytes.Buffer)
	if err := pc.pngEncoder.Encode(buff, combinedImg); err != nil {
		return nil, fmt.Errorf("Error Encoding Image: %w", err)
	}

	return buff, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
