package image

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

func HexToRGBA(hexColor string) (color.RGBA, error) {
	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hexColor)
	}

	r, err := strconv.ParseUint(hexColor[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseUint(hexColor[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseUint(hexColor[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

// Calculate Euclidan Distance between the two pixles
func MeasureColorEuclideanDist(c1, c2 color.RGBA) uint32 {
	dr := int32(c1.R) - int32(c2.R)
	dg := int32(c1.G) - int32(c2.G)
	db := int32(c1.B) - int32(c2.B)
	return uint32(dr*dr + dg*dg + db*db)
}

// Compare how similar two pixels are
func isSimilar(c1, c2 color.RGBA, threshold uint32) bool {
	return MeasureColorEuclideanDist(c1, c2) < threshold
}
