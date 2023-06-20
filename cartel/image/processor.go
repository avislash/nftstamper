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

type pledgeHands struct {
	hands map[string]map[string]map[string]image.Image
}

func (p *pledgeHands) getMAYCHands() map[string]map[string]image.Image {
	hands, _ := p.hands["mayc"]
	return hands
}

func (p *pledgeHands) getMutantHoundsHands() map[string]map[string]image.Image {
	hands, _ := p.hands["mutant_hounds"]
	return hands
}

type Processor struct {
	image.Combiner
	bowls       map[string]image.Image //map of backgrounds to bowls
	pledgeHands pledgeHands            //map[string]map[string]image.Image //map of traits to colorss
	nfdMerch    Merch
	nfdSuit     image.Image
	apeBags     map[string]image.Image
}

func NewProcessor(config config.ImageProcessorConfig) (*Processor, error) {
	bowls, err := buildImageMap(config.GMMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building GM Bowl Image Mappings: %w", err)
	}

	pledgeHands := pledgeHands{hands: make(map[string]map[string]map[string]image.Image)}
	hands := make(map[string]map[string]image.Image)
	for trait, mappings := range config.PledgeHands.MAYC {
		colorMap, err := buildImageMap(mappings.Colors)
		if err != nil {
			return nil, fmt.Errorf("Error building Hand Image Mappings for %s: %w", trait, err)
		}
		hands[trait] = colorMap
	}
	pledgeHands.hands["mayc"] = hands

	hands = make(map[string]map[string]image.Image)
	for trait, mappings := range config.PledgeHands.Hounds {
		colorMap, err := buildImageMap(mappings.Colors)
		if err != nil {
			return nil, fmt.Errorf("Error building Hand Image Mappings for %s: %w", trait, err)
		}
		hands[trait] = colorMap
	}
	pledgeHands.hands["mutant_hounds"] = hands

	apeBags, err := buildImageMap(config.ApeBagMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building Ape Bag Mappings: %w", err)
	}

	nfdMerch, err := buildMerch(config.NFDMerchMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building NFD Merch Image Mappings: %w", err)
	}

	suit, err := getImageFromFile(config.Suit)
	if err != nil {
		return nil, fmt.Errorf("Error loading Suit: %w", err)
	}

	return &Processor{
		//Combined Hound images are too big to process and return to discord before timing out
		Combiner:    image.NewPNGCombiner(image.WithBestSpeedPNGCompression()),
		bowls:       bowls,
		pledgeHands: pledgeHands,
		nfdMerch:    nfdMerch,
		nfdSuit:     suit,
		apeBags:     apeBags,
	}, nil
}

func (p *Processor) OverlayBowl(hound image.Image, background string) (*bytes.Buffer, error) {
	bowl, exists := p.bowls[background]
	if !exists {
		return nil, fmt.Errorf("No bowl file found for background: %s", background)
	}
	return p.EncodeImage(p.CombineImages(hound, bowl))
}

func (p *Processor) OverlayNFDSuit(ape image.Image) (*bytes.Buffer, error) {
	return p.EncodeImage(p.CombineImages(ape, p.nfdSuit))
}

func (p *Processor) OverlayHandMAYC(ape image.Image, metadata metadata.MAYCMetadata, color string) (*bytes.Buffer, error) {
	key := "default"
	maycHands := p.pledgeHands.getMAYCHands()
	hands := maycHands[key]

	if override, exists := maycHands[metadata.Clothes]; exists {
		key = metadata.Clothes
		hands = override
	}

	if len(hands) == 0 {
		return nil, fmt.Errorf("No Default or Trait Color Map Defined at key: %s", key)
	}

	hand, exists := hands[color]
	if !exists {
		return nil, fmt.Errorf("No hand image found for %s", color)
	}

	return p.EncodeImage(p.CombineImages(ape, hand))
}

func (p *Processor) OverlayHandHound(hound image.Image, _ metadata.HoundMetadata, color string) (*bytes.Buffer, error) {
	key := "default"
	houndHands := p.pledgeHands.getMutantHoundsHands()
	hands := houndHands[key]

	if len(hands) == 0 {
		return nil, fmt.Errorf("No Default or Trait Color Map Defined at key: %s", key)
	}

	hand, exists := hands[color]
	if !exists {
		return nil, fmt.Errorf("No hand image found for %s", color)
	}

	return p.EncodeImage(p.CombineImages(hound, hand))
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

func (p *Processor) OverlayApeBag(ape image.Image, metadata metadata.MAYCMetadata) (*bytes.Buffer, error) {
	bag := p.apeBags["default"]
	if override, exists := p.apeBags[metadata.Mouth]; exists {
		bag = override
	}
	return p.EncodeImage(p.CombineImages(ape, bag))
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
