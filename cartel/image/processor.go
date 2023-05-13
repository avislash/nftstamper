package image

import (
	"bytes"
	"fmt"
	"os"

	"github.com/avislash/nftstamper/cartel/config"
	"github.com/avislash/nftstamper/lib/image"
)

type Processor struct {
	image.Combiner
	bowls map[string]image.Image //map of backgrounds to bowls
}

func NewProcessor(config config.ImageProcessorConfig) (*Processor, error) {
	decoder := &image.PNGDecoder{}
	bowls := make(map[string]image.Image)

	for background, bowlFile := range config.GMMappings {
		file, err := os.Open(bowlFile)
		if err != nil {
			return nil, fmt.Errorf("Unable to open %s: %w", bowlFile, err)
		}
		defer file.Close()

		img, err := decoder.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode image file %s: %w", bowlFile, err)
		}
		bowls[background] = img
	}
	return &Processor{
		//Combined Hound images are too big to process and return to discord before timing out
		Combiner: image.NewPNGCombiner(image.WithBestSpeedPNGCompression()),
		bowls:    bowls,
	}, nil
}

func (p *Processor) OverlayBowl(hound image.Image, background string) (*bytes.Buffer, error) {
	bowl, exists := p.bowls[background]
	if !exists {
		return nil, fmt.Errorf("No bowl file found for background: %s", background)
	}
	return p.EncodeImage(p.CombineImages(hound, bowl))
}
