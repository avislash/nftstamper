package image

import (
	"bytes"
	"fmt"
	"os"

	"github.com/avislash/nftstamper/ape/config"
	"github.com/avislash/nftstamper/ape/metadata"
	"github.com/avislash/nftstamper/lib/image"
	"github.com/avislash/nftstamper/lib/log"
)

type Processor struct {
	image.Combiner
	logger            *log.SugaredLogger
	mugs              map[string]image.Image //map of base armors to mug images
	baseGMSmoke       image.Image
	azulGMSmoke       image.Image //Azuls are different size than base armors
	baseGMSmokeBorder image.Image
	azulGMSmokeBorder image.Image //Azuls are different size than base armors
	opacityFilter     *OpacityFilter
}

func NewProcessor(config config.ImageProcessorConfig, logger *log.SugaredLogger) (*Processor, error) {
	mugs := make(map[string]image.Image)

	baseGMSmoke, err := getImageFromFile(config.BaseGMSmoke)
	if err != nil {
		return nil, err
	}

	baseGMSmokeBorder, err := getImageFromFile(config.BaseGMSmokeBorder)
	if err != nil {
		return nil, err
	}

	azulGMSmoke, err := getImageFromFile(config.AzulGMSmoke)
	if err != nil {
		return nil, err
	}

	azulGMSmokeBorder, err := getImageFromFile(config.AzulGMSmokeBorder)
	if err != nil {
		return nil, err
	}

	for baseArmor, mugFile := range config.GMMappings {
		img, err := getImageFromFile(mugFile)
		if err != nil {
			return nil, err
		}
		mugs[baseArmor] = img
	}

	return &Processor{
		Combiner:          image.NewPNGCombiner(),
		logger:            logger,
		mugs:              mugs,
		baseGMSmoke:       baseGMSmoke,
		baseGMSmokeBorder: baseGMSmokeBorder,
		azulGMSmoke:       azulGMSmoke,
		azulGMSmokeBorder: azulGMSmokeBorder,
		opacityFilter:     NewOpacityFilter(config.Filters.Opacity),
	}, nil
}

func (p *Processor) OverlayMug(sentinel image.Image, metadata metadata.SentinelMetadata) (*bytes.Buffer, error) {
	sentinelHand, exists := p.mugs[metadata.BaseArmor]
	if !exists {
		return nil, fmt.Errorf("No mug file found for base armor: %s", metadata.BaseArmor)
	}
	gmSmoke := p.baseGMSmoke
	gmSmokeBorder := p.baseGMSmokeBorder
	if metadata.BaseArmor == "azul" {
		gmSmoke = p.azulGMSmoke
		gmSmokeBorder = p.azulGMSmokeBorder
	}

	filters := p.opacityFilter.Filters[metadata.BaseArmor]
	opacity := filters.Default
	if len(metadata.Body) != 0 {
		op, exists := filters.Weights[metadata.Body]
		if exists {
			opacity = op
		}
	}
	p.logger.Debugf("Opacity filter for Sentinel with metadata %+v set to: %f", metadata.Attributes, opacity)
	gmSmoke = p.AdjustImageOpacity(gmSmoke, opacity)
	smokeWithBorder := p.CombineImages(gmSmoke, gmSmokeBorder)
	coffeMugWithSmoke := p.CombineImages(sentinelHand, smokeWithBorder)
	return p.EncodeImage(p.CombineImages(sentinel, coffeMugWithSmoke))
}

func getImageFromFile(filename string) (image.Image, error) {
	decoder := &image.PNGDecoder{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %s: %w", filename, err)
	}
	defer file.Close()

	img, err := decoder.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode image file %s: %w", filename, err)
	}
	return img, err
}
