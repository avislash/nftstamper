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
