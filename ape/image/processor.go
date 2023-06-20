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

type mug struct {
	path image.Image
	lab  image.Image
}

type smoke struct {
	lab    image.Image
	path   image.Image
	border image.Image
}

type Processor struct {
	image.Combiner
	logger                *log.SugaredLogger
	hbdMappings           map[string]image.Image
	mugs                  map[string]mug //map of base armors to mug images
	baseGMSmokeProperties smoke
	azulGMSmokeProperties smoke //Azuls are different size than base armors
	opacityFilter         *OpacityFilter
}

func NewProcessor(config config.ImageProcessorConfig, logger *log.SugaredLogger) (*Processor, error) {

	smoke, err := buildSmoke(config)
	if err != nil {
		return nil, err
	}

	mugs, err := buildMugs(config)
	if err != nil {
		return nil, err
	}

	hbdMappings, err := buildImageMap(config.HBDMappings)
	if err != nil {
		return nil, err
	}

	return &Processor{
		Combiner:              image.NewPNGCombiner(),
		logger:                logger,
		hbdMappings:           hbdMappings,
		mugs:                  mugs,
		baseGMSmokeProperties: smoke["base"],
		azulGMSmokeProperties: smoke["azul"],
		opacityFilter:         NewOpacityFilter(config.Filters.Opacity),
	}, nil
}

func (p *Processor) OverlayMug(sentinel image.Image, metadata metadata.SentinelMetadata) (*bytes.Buffer, error) {

	mug, found := p.getMug(metadata)
	if !found {
		return nil, fmt.Errorf("No mug found for base armor: %s", metadata.BaseArmor)
	}

	smoke, border := p.getSmoke(metadata)

	if smoke == nil {
		err := fmt.Errorf("No smoke defined")
		p.logger.Error(err.Error())
		return nil, err
	}

	if border == nil {
		err := fmt.Errorf("No border defined")
		p.logger.Error(err.Error())
		return nil, err
	}

	smoke = p.AdjustSmokeOpacity(smoke, metadata)
	smokeWithBorder := p.CombineImages(smoke, border)
	coffeMugWithSmoke := p.CombineImages(mug, smokeWithBorder)
	return p.EncodeImage(p.CombineImages(sentinel, coffeMugWithSmoke))
}

func (p *Processor) OverlayHBD(sentinel image.Image, metadata metadata.SentinelMetadata) (*bytes.Buffer, error) {
	hbdHand, exists := p.hbdMappings[metadata.BaseArmor]
	if !exists {
		return nil, fmt.Errorf("No HBD Image found for base armor: %s", metadata.BaseArmor)
	}
	return p.EncodeImage(p.CombineImages(sentinel, hbdHand))
}

func (p *Processor) AdjustSmokeOpacity(smoke image.Image, metadata metadata.SentinelMetadata) image.Image {
	filters := p.opacityFilter.Filters[metadata.BaseArmor]
	opacity := filters.Default

	if len(metadata.Body) != 0 {
		op, exists := filters.Weights[metadata.Body]
		if exists {
			opacity = op
		}
	}

	if len(metadata.Head) != 0 {
		op, exists := filters.Weights[metadata.Head]
		if exists {
			opacity = op
		}
	}

	p.logger.Debugf("Opacity filter for Sentinel with metadata %+v set to: %f", metadata.Attributes, opacity)
	return p.AdjustImageOpacity(smoke, opacity)
}

func (p *Processor) getMug(metadata metadata.SentinelMetadata) (image.Image, bool) {
	mug, exists := p.mugs[metadata.Attributes.BaseArmor]
	if !exists {
		return nil, false
	}

	if metadata.Attributes.Head == "path robe" {
		return mug.path, true
	}
	return mug.lab, true
}

func (p *Processor) getSmoke(metadata metadata.SentinelMetadata) (image.Image, image.Image) {
	if metadata.BaseArmor == "azul" {
		return p.azulGMSmokeProperties.lab, p.azulGMSmokeProperties.border
	}

	if metadata.Attributes.Head == "path robe" {
		return p.baseGMSmokeProperties.path, p.baseGMSmokeProperties.border
	}
	return p.baseGMSmokeProperties.lab, p.baseGMSmokeProperties.border

}

func buildMugs(config config.ImageProcessorConfig) (map[string]mug, error) {
	mugs := make(map[string]mug)
	for baseArmor, mugFiles := range config.GMMappings {
		mug := mug{}
		for name, file := range mugFiles {
			img, err := getImageFromFile(file)
			if err != nil {
				return nil, err
			}

			switch name {
			case "path":
				mug.path = img
			case "lab":
				mug.lab = img
			}
		}
		mugs[baseArmor] = mug
	}
	return mugs, nil
}
func buildSmoke(config config.ImageProcessorConfig) (map[string]smoke, error) {
	baseGMSmoke, err := getImageFromFile(config.BaseGMSmoke)
	if err != nil {
		return nil, err
	}

	baseGMSmokeBorder, err := getImageFromFile(config.BaseGMSmokeBorder)
	if err != nil {
		return nil, err
	}

	pathGMSmoke, err := getImageFromFile(config.BaseGMSmokePath)
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

	return map[string]smoke{
		"base": smoke{baseGMSmoke, pathGMSmoke, baseGMSmokeBorder},
		"azul": smoke{azulGMSmoke, nil, azulGMSmokeBorder},
	}, nil
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

func buildImageMap(mappings map[string]string) (map[string]image.Image, error) {
	imageMap := make(map[string]image.Image)
	for key, imgFile := range mappings {
		img, err := getImageFromFile(imgFile)
		if err != nil {
			return nil, fmt.Errorf("Error building image for %s: %w", key, err)
		}
		imageMap[key] = img
	}
	return imageMap, nil
}
