package image

import (
	"bytes"
	"fmt"
	"os"

	"github.com/avislash/nftstamper/config"
	"github.com/avislash/nftstamper/lib/image"
)

type Processor struct {
	image.Combiner
	mugs map[string]image.Image //map of base armors to mug images
}

func NewProcessor(config config.ImageProcessorConfig) (*Processor, error) {
	decoder := &image.PNGDecoder{}
	mugs := make(map[string]image.Image)

	for baseArmor, mugFile := range config.GMMappings {
		file, err := os.Open(mugFile)
		if err != nil {
			return nil, fmt.Errorf("Unable to open %s: %w", mugFile, err)
		}
		defer file.Close()

		img, err := decoder.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode image file %s: %w", mugFile, err)
		}
		mugs[baseArmor] = img
	}
	return &Processor{
		Combiner: image.NewPNGCombiner(),
		mugs:     mugs}, nil
}

func (p *Processor) OverlayMug(sentinel image.Image, baseArmor string) (*bytes.Buffer, error) {
	sentinelHand, exists := p.mugs[baseArmor]
	if !exists {
		return nil, fmt.Errorf("No mug file found for base armor: %s", baseArmor)
	}
	return p.EncodeImage(p.CombineImages(sentinel, sentinelHand))
}
