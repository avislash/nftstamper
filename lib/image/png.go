package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

var _ Combiner = (*PNGCombiner)(nil)

type PNGCombinerOption func(p *PNGCombiner)

type PNGCombiner struct {
	pngEncoder png.Encoder
}

func WithDefaultPNGCompression() PNGCombinerOption {
	return WithCompressionLevel(png.DefaultCompression)
}

func WithNoPNGCompression() PNGCombinerOption {
	return WithCompressionLevel(png.NoCompression)
}

func WithBestSpeedPNGCompression() PNGCombinerOption {
	return WithCompressionLevel(png.BestSpeed)
}

func WithBestPNGCompression() PNGCombinerOption {
	return WithCompressionLevel(png.BestCompression)
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

func (pc *PNGCombiner) CombineImages(img1, img2 Image) Image {
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

	return combinedImg
}

// Adjusts the image opacity of non-transparent pixels to the specified opacity
// The opacity adjustment is made using he over-composition mode of the Porter-Duff algorithm.
// Using over-composition since this (based on observation) allows for the best
// result when blending foreground over background in the current use case.
// The other modes can be added later and this function can be refactored if needed.
func (pc *PNGCombiner) AdjustImageOpacity(img Image, opacity float64) Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	// Map opacity (0-1) to pixel value ranging from 0-255
	alpha := uint8(opacity * 0xFF)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, originalAlpha := rgba.At(x, y).RGBA()

			// Calculate the new alpha value while preserving transparency
			newAlpha := uint32(originalAlpha>>8) * uint32(alpha) / 0xFF

			if originalAlpha == 0 {
				// Transparent pixel, no adjustment needed
				continue
			}

			// Calculate the adjusted RGB values based on the new alpha value
			adjustedR := uint8(uint32(r>>8) * newAlpha / uint32(originalAlpha>>8))
			adjustedG := uint8(uint32(g>>8) * newAlpha / uint32(originalAlpha>>8))
			adjustedB := uint8(uint32(b>>8) * newAlpha / uint32(originalAlpha>>8))

			// Update the pixel with the adjusted alpha and RGB values
			rgba.SetRGBA(x, y, color.RGBA{
				R: adjustedR,
				G: adjustedG,
				B: adjustedB,
				A: uint8(newAlpha),
			})
		}
	}

	return rgba
}

func (pc *PNGCombiner) HexChromaKeySwap(img Image, chromaKey, newColor string) (Image, error) {
	_chromaKey, err := HexToRGBA(chromaKey)
	if err != nil {
		return nil, fmt.Errorf("Error converting chromaKey %s to RGBA: %w", chromaKey, err)
	}

	newKey, err := HexToRGBA(newColor)
	if err != nil {
		return nil, fmt.Errorf("Error converting newColor %s to RGBA: %w", newColor, err)
	}

	dstImg := image.NewRGBA(img.Bounds())
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// Get the color of the pixel
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			// If the pixel color is the chroma key color, set it to bg color
			// If it's not, copy it to the destination image
			if c.R == _chromaKey.R && c.G == _chromaKey.G && c.B == _chromaKey.B {
				dstImg.Set(x, y, newKey)
			} else {
				dstImg.Set(x, y, c)
			}
		}
	}

	return dstImg, nil
}

func (pc *PNGCombiner) FilterOutBackgroundColor(img Image, bgKey string, threshold uint32) (Image, error) {
	bgColor, err := HexToRGBA(bgKey)
	if err != nil {
		return nil, fmt.Errorf("Error converting background key %s to RGBA: %w", bgKey, err)
	}

	bounds := img.Bounds()
	subject := image.NewRGBA(bounds)

	// Step 2: Loop over each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			r, g, b, a := img.At(x, y).RGBA()
			pixelColor := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}

			if isSimilar(pixelColor, bgColor, threshold) {
				subject.Set(x, y, color.Transparent)
			} else {
				subject.Set(x, y, pixelColor)
			}
		}
	}

	return subject, nil
}

func (pc *PNGCombiner) EncodeImage(img Image) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	if err := pc.pngEncoder.Encode(buff, img); err != nil {
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
