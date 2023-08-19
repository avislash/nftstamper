package image

import (
	"bytes"
	"fmt"
	"os"
	"strings"

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

type maskMappings struct {
	chromaKey     string
	defaultMask   image.Image
	traitMappings map[string]image.Image
}

type suitMappings struct {
	masks    map[string]maskMappings
	skipMask map[string]bool
	suits    map[string]map[string]image.Image
}

type coffeeMugMappings struct {
	liquids map[string]image.Image
	steam   map[string]image.Image
	logos   map[string]image.Image
	furs    map[string]image.Image
}

type houndTraitMappings struct {
	faces  map[string]image.Image
	forms  map[string]image.Image
	heads  map[string]image.Image
	legs   map[string]image.Image
	mouths map[string]image.Image
	noses  map[string]image.Image
	torsos map[string]image.Image
}

type Processor struct {
	image.Combiner
	bowls                   map[string]image.Image //map of backgrounds to bowls
	pledgeHands             pledgeHands            //map of traits to colors
	nfdMerch                Merch
	suitMappings            suitMappings
	apeBags                 map[string]image.Image
	maycBackgroundColorKeys map[string]string
	maycCoffeeMugMappings   coffeeMugMappings
	nflJerseyMappings       map[string]map[string]image.Image
	houndTraitMappings      houndTraitMappings
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

	suitMappings, err := buildSuitMappings(config.SuitMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building Suit Mappings: %w", err)
	}

	maycCoffeeMugMappings, err := buildCoffeeMugMappings(config.MAYCCoffeeMugMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building MAYC Coffee Mug Mappings: %w", err)
	}

	jerseyMappings := make(map[string]map[string]image.Image)
	for key, mappings := range config.NFLJerseyMappings {
		imgMap, err := buildImageMap(mappings)
		if err != nil {
			return nil, fmt.Errorf("Error building NFL Jersey Mappings for %s: %w", key, err)
		}
		jerseyMappings[key] = imgMap
	}

	houndTraitMappings, err := buildHoundTraitMappings(config.HoundTraitMappings)
	if err != nil {
		return nil, fmt.Errorf("Error building Hound Trait Mappings: %w", err)
	}

	return &Processor{
		//Combined Hound images are too big to process and return to discord before timing out
		Combiner:                image.NewPNGCombiner(image.WithBestSpeedPNGCompression()),
		bowls:                   bowls,
		pledgeHands:             pledgeHands,
		nfdMerch:                nfdMerch,
		suitMappings:            suitMappings,
		apeBags:                 apeBags,
		maycBackgroundColorKeys: config.MAYCBackgroundColorKeys,
		maycCoffeeMugMappings:   maycCoffeeMugMappings,
		nflJerseyMappings:       jerseyMappings,
		houndTraitMappings:      houndTraitMappings,
	}, nil
}

func (p *Processor) OverlayBowl(hound image.Image, background string) (*bytes.Buffer, error) {
	bowl, exists := p.bowls[background]
	if !exists {
		return nil, fmt.Errorf("No bowl file found for background: %s", background)
	}
	return p.EncodeImage(p.CombineImages(hound, bowl))
}

func (p *Processor) OverlaySuit(suit string, ape image.Image, metadata metadata.MAYCMetadata) (*bytes.Buffer, error) {
	if strings.Contains(metadata.Name, "mega") {
		return nil, fmt.Errorf("Mega Mutants not supported")
	}

	skipMask, _ := p.suitMappings.skipMask[metadata.Clothes]
	if !skipMask {
		background := metadata.Background[3:]
		bgKey, exists := p.maycBackgroundColorKeys[background]
		if !exists {
			return nil, fmt.Errorf("No background key defined for %s", background)
		}

		mask, chromaKey := p.getSuitMask(metadata)

		mask, err := p.HexChromaKeySwap(mask, chromaKey, bgKey)
		if err != nil {
			return nil, fmt.Errorf("Error Chroma Keying Mask: %w", err)
		}
		ape = p.CombineImages(ape, mask)
	}

	var suits map[string]image.Image
	switch {
	case strings.Contains(metadata.Mouth, "m2 bored"):
		suits = p.suitMappings.suits["m2 bored"]
	default:
		suits = p.suitMappings.suits["default"]
	}

	suitImg, exists := suits[suit]
	if !exists {
		return nil, fmt.Errorf("No suit loaded for %s", suit)
	}

	return p.EncodeImage(p.CombineImages(ape, suitImg))
}

func (p *Processor) OverlayCoffeeMug(ape image.Image, metadata metadata.MAYCMetadata, liquid, logo string) (*bytes.Buffer, error) {
	mug, exists := p.maycCoffeeMugMappings.furs[metadata.Fur]
	if !exists {
		return nil, fmt.Errorf("No Coffee Mug found for fur: %s", metadata.Fur)
	}

	liquidImg, exists := p.maycCoffeeMugMappings.liquids[liquid]
	if !exists {
		return nil, fmt.Errorf("No Liquid found for liquid: %s", liquid)
	}
	steam, _ := p.maycCoffeeMugMappings.steam[liquid] //Steam is optional and not necessary for final image

	logoImg, exists := p.maycCoffeeMugMappings.logos[logo]
	if !exists {
		return nil, fmt.Errorf("No Logo found for logo: %s", logo)
	}

	//Mug goes over liquid layer whilst everything goes over the mug
	mug = p.CombineImages(liquidImg, mug)
	if steam != nil {
		mug = p.CombineImages(mug, steam)
	}
	mug = p.CombineImages(mug, logoImg)
	return p.EncodeImage(p.CombineImages(ape, mug))
}

func (p *Processor) getSuitMask(metadata metadata.MAYCMetadata) (image.Image, string) {
	if strings.Contains(metadata.Mouth, "m1 bored") {
		maskMappings := p.suitMappings.masks["m1 bored"]
		return getMaskFromMetadata(maskMappings, metadata), maskMappings.chromaKey
	}

	if strings.Contains(metadata.Mouth, "m2 bored") {
		maskMappings := p.suitMappings.masks["m2 bored"]
		return getMaskFromMetadata(maskMappings, metadata), maskMappings.chromaKey
	}

	if maskMappings, exists := p.suitMappings.masks[metadata.Mouth]; exists {
		return getMaskFromMetadata(maskMappings, metadata), maskMappings.chromaKey
	}

	return getMaskFromMetadata(p.suitMappings.masks["default"], metadata), p.suitMappings.masks["default"].chromaKey
}

func getMaskFromMetadata(maskMappings maskMappings, metadata metadata.MAYCMetadata) image.Image {
	traitMappings := maskMappings.traitMappings

	if mask, exists := traitMappings[metadata.Mouth]; exists {
		return mask
	}

	if mask, exists := traitMappings[metadata.Hat]; exists {
		return mask
	}

	if mask, exists := traitMappings[metadata.Clothes]; exists {
		return mask
	}

	return maskMappings.defaultMask
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

func (p *Processor) OverlayHoundJersey(hound image.Image, metadata metadata.HoundMetadata, team string) (*bytes.Buffer, error) {
	jerseys, exists := p.nflJerseyMappings["default"]
	if !exists {
		return nil, fmt.Errorf("No default Jerseys loaded")
	}

	jersey, exists := jerseys[team]
	if !exists {
		return nil, fmt.Errorf("No jersey loaded for %s", team)
	}

	return p.EncodeImage(p.CombineImages(hound, jersey))
}

func (p *Processor) CutoutHound(metadata metadata.HoundMetadata) (*bytes.Buffer, error) {
	hound, err := p.generateHound(metadata)
	if err != nil {
		return nil, err
	}
	return p.EncodeImage(hound)
}

func (p *Processor) generateHound(metadata metadata.HoundMetadata) (image.Image, error) {
	var hound image.Image
	traitMap := p.houndTraitMappings

	hound, exists := traitMap.forms[metadata.Form]
	if !exists {
		return hound, fmt.Errorf("No form loaded for: %s", metadata.Form)
	}

	if leg, exists := traitMap.legs[metadata.Leg]; exists {
		hound = p.CombineImages(hound, leg)
	}

	if torso, exists := traitMap.torsos[metadata.Torso]; exists {
		hound = p.CombineImages(hound, torso)
	}

	if face, exists := traitMap.faces[metadata.Face]; exists {
		hound = p.CombineImages(hound, face)
	}

	if mouth, exists := traitMap.mouths[metadata.Mouth]; exists {
		hound = p.CombineImages(hound, mouth)
	}

	if head, exists := traitMap.heads[metadata.Head]; exists {
		hound = p.CombineImages(hound, head)
	}

	if nose, exists := traitMap.noses[metadata.Nose]; exists {
		hound = p.CombineImages(hound, nose)
	}

	return hound, nil
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

func buildSuitMappings(suitConfig config.SuitMappings) (suitMappings, error) {
	masks := make(map[string]maskMappings)
	for key, maskMapping := range suitConfig.Masks {
		defaultMask, err := getImageFromFile(maskMapping.Default)
		if err != nil {
			return suitMappings{}, fmt.Errorf("Error loading default mask for %s: %w", key, err)
		}

		traitMappings, err := buildImageMap(maskMapping.TraitMappings)
		if err != nil {
			return suitMappings{}, fmt.Errorf("Error build mask trait mappings for %s: %w", key, err)
		}

		masks[key] = maskMappings{
			defaultMask:   defaultMask,
			traitMappings: traitMappings,
			chromaKey:     maskMapping.ChromaKey,
		}

	}

	suits := make(map[string]map[string]image.Image)
	for key, suitMap := range suitConfig.Suits {
		suitImgMap, err := buildImageMap(suitMap)
		if err != nil {
			return suitMappings{}, fmt.Errorf("Error building suit image map for %s: %w", key, err)
		}
		suits[key] = suitImgMap
	}

	skipMask := make(map[string]bool)
	for _, clothing := range suitConfig.SkipMask {
		skipMask[clothing] = true
	}

	return suitMappings{
		masks:    masks,
		skipMask: skipMask,
		suits:    suits,
	}, nil
}

func buildCoffeeMugMappings(coffeeMugConfig config.CoffeeMugMappings) (coffeeMugMappings, error) {
	furs, err := buildImageMap(coffeeMugConfig.Furs)
	if err != nil {
		return coffeeMugMappings{}, fmt.Errorf("Error building fur image map: %w", err)
	}

	logos, err := buildImageMap(coffeeMugConfig.Logos)
	if err != nil {
		return coffeeMugMappings{}, fmt.Errorf("Error building logo image map: %w", err)
	}

	liquids, err := buildImageMap(coffeeMugConfig.Liquids)
	if err != nil {
		return coffeeMugMappings{}, fmt.Errorf("Error building liquids image map: %w", err)
	}

	steam, err := buildImageMap(coffeeMugConfig.Steam)
	if err != nil {
		return coffeeMugMappings{}, fmt.Errorf("Error building steam image map: %w", err)
	}

	return coffeeMugMappings{
		liquids: liquids,
		steam:   steam,
		logos:   logos,
		furs:    furs,
	}, nil

}

func buildHoundTraitMappings(houndTraitConfig config.HoundTraitMappings) (houndTraitMappings, error) {
	faces, err := buildImageMap(houndTraitConfig.Faces)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building face image map: %w", err)
	}

	forms, err := buildImageMap(houndTraitConfig.Forms)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building form image map: %w", err)
	}

	heads, err := buildImageMap(houndTraitConfig.Heads)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building head image map: %w", err)
	}

	legs, err := buildImageMap(houndTraitConfig.Legs)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building leg image map: %w", err)
	}

	mouths, err := buildImageMap(houndTraitConfig.Mouths)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building mouth image map: %w", err)
	}

	noses, err := buildImageMap(houndTraitConfig.Noses)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building nose image map: %w", err)
	}

	torsos, err := buildImageMap(houndTraitConfig.Torsos)
	if err != nil {
		return houndTraitMappings{}, fmt.Errorf("Error building torso image map: %w", err)
	}

	return houndTraitMappings{
		faces:  faces,
		forms:  forms,
		heads:  heads,
		legs:   legs,
		mouths: mouths,
		noses:  noses,
		torsos: torsos,
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
