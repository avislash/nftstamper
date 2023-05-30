package image

import (
	"bytes"
	"fmt"
	"os"

	"github.com/avislash/nftstamper/cartel/config"
	"github.com/avislash/nftstamper/cartel/metadata"
	"github.com/avislash/nftstamper/lib/image"
)

type Merch struct {
	Default     image.Image
	XL          image.Image
	Flame       image.Image
	XLFlame     image.Image
	Hats        map[string]image.Image
	FlameTraits map[string]string
	XLTraits    map[string]string
}
type Processor struct {
	image.Combiner
	bowls    map[string]image.Image //map of backgrounds to bowls
	nfdMerch Merch
}

func NewProcessor(config config.ImageProcessorConfig) (*Processor, error) {
	bowls, err := buildImageMap(config.GMMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building GM Bowl Image Mappings: %w", err)
	}

	nfdMerch, err := buildMerch(config.NFDMerchMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building NFD Merch Image Mappings: %w", err)
	}

	return &Processor{
		//Combined Hound images are too big to process and return to discord before timing out
		Combiner: image.NewPNGCombiner(image.WithBestSpeedPNGCompression()),
		bowls:    bowls,
		nfdMerch: nfdMerch,
	}, nil
}

func (p *Processor) OverlayBowl(hound image.Image, background string) (*bytes.Buffer, error) {
	bowl, exists := p.bowls[background]
	if !exists {
		return nil, fmt.Errorf("No bowl file found for background: %s", background)
	}
	return p.EncodeImage(p.CombineImages(hound, bowl))
}

func (p *Processor) OverlayNFDMerch(hound image.Image, metadata metadata.HoundMetadata) (*bytes.Buffer, error) {
	var torso image.Image
	_, isXL := p.nfdMerch.XLTraits[metadata.Torso]
	//Check to see if we need flame shirt
	if _, isFlame := p.nfdMerch.FlameTraits[metadata.Mouth]; isFlame {
		torso = p.nfdMerch.Flame
		if isXL {
			torso = p.nfdMerch.XLFlame
		}
	} else {
		torso = p.nfdMerch.Default
		if isXL {
			torso = p.nfdMerch.XL
		}
	}

	var merch image.Image = torso
	if hat, exists := p.nfdMerch.Hats[metadata.Face]; exists {
		merch = p.CombineImages(merch, hat)
	}

	return p.EncodeImage(p.CombineImages(hound, merch))
}

func buildImageMap(imageFiles map[string]string) (map[string]image.Image, error) {
	mappings := make(map[string]image.Image)
	for trait, imageFile := range imageFiles {
		img, err := getImageFromFile(imageFile)
		if err != nil {
			return nil, err
		}

		mappings[trait] = img
	}
	return mappings, nil
}

func buildMerch(merchConfig config.MerchMappings) (Merch, error) {
	_default, err := getImageFromFile(merchConfig.Default)
	if err != nil {
		return Merch{}, fmt.Errorf("Error loading default image: %w", err)
	}

	flame, err := getImageFromFile(merchConfig.Flame)
	if err != nil {
		return Merch{}, fmt.Errorf("Error loading flame image: %w", err)
	}

	xl, err := getImageFromFile(merchConfig.XL)
	if err != nil {
		return Merch{}, fmt.Errorf("Error loading XL image: %w", err)
	}

	xlFlame, err := getImageFromFile(merchConfig.XLFlame)
	if err != nil {
		return Merch{}, fmt.Errorf("Error loading XL Flame image: %w", err)
	}

	hats, err := buildImageMap(merchConfig.Hats)
	if err != nil {
		return Merch{}, fmt.Errorf("Error loading hat images: %w", err)
	}

	return Merch{
		Default:     _default,
		Flame:       flame,
		XL:          xl,
		XLFlame:     xlFlame,
		Hats:        hats,
		FlameTraits: merchConfig.FlameTraits,
		XLTraits:    merchConfig.XLTraits,
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
