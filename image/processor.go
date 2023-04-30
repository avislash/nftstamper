package image

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Processor struct {
	mugs map[string]image.Image
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) OverlayMug(sentinel image.Image, baseArmor string) (*bytes.Buffer, error) {
	// Open the second image file
	sentinelHandFile, err := os.Open("./mugs/trippyMug.png")
	if err != nil {
		panic(err)
	}
	defer sentinelHandFile.Close()

	sentinelHand, err := png.Decode(sentinelHandFile)
	if err != nil {
		panic(err)
	}

	// Create a new image with the size of the larger image
	combinedWidth := max(sentinel.Bounds().Max.X, sentinelHand.Bounds().Max.X)
	combinedHeight := max(sentinel.Bounds().Max.Y, sentinelHand.Bounds().Max.Y)
	combinedImg := image.NewRGBA(image.Rect(0, 0, combinedWidth, combinedHeight))

	// Draw the first image onto the combined image
	draw.Draw(combinedImg, sentinel.Bounds(), sentinel, image.ZP, draw.Src)

	// Draw the second image onto the combined image with an offset
	offset := image.Pt((combinedWidth-sentinelHand.Bounds().Dx())/2, (combinedHeight-sentinelHand.Bounds().Dy())/2)
	drawRect := sentinelHand.Bounds()
	drawRect = drawRect.Add(offset)
	drawRect = drawRect.Intersect(combinedImg.Bounds())
	drawRect = drawRect.Sub(offset)
	drawRect = drawRect.Add(offset)
	drawRect = drawRect.Intersect(sentinelHand.Bounds())
	drawRect = drawRect.Sub(offset)
	draw.Draw(combinedImg, drawRect, sentinelHand, sentinelHand.Bounds().Min, draw.Over)

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, combinedImg); err != nil {
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
