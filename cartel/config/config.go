package config

import (
	"fmt"

	"github.com/avislash/nftstamper/lib/config"
)

type Config struct {
	HoundsMetadataEndpoint string               `yaml:"hounds_metadata_endpoint"`
	MAYCMetadataEndpoint   string               `yaml:"mayc_metadata_endpoint"`
	BAYCMetadataEndpoint   string               `yaml:"bayc_metadata_endpoint"`
	BotToken               string               `yaml:"discord_bot_token"`
	IPFSEndpoint           string               `yaml:"ipfs_endpoint"`
	ImageProcessorConfig   ImageProcessorConfig `yaml:"image_processor_mappings"`
}

type HandMapping struct {
	Hand map[string]string
}

type MerchMappings struct {
	Default     string            `yaml:"default"`
	Flame       string            `yaml:"default_flame"`
	XL          string            `yaml:"xl"`
	XLFlame     string            `yaml:"xl_flame"`
	FlameTraits map[string]string `yaml:"flame_traits"`
	XLTraits    map[string]string `yaml:"xl_traits"`
	Hats        map[string]string `yaml:"hats"`
}

type HandMappings struct {
	Colors map[string]string `yaml:"colors"`
}

type PledgeMappings struct {
	Hounds map[string]HandMappings `yaml:"mutant_hounds"`
	MAYC   map[string]HandMappings `yaml:"mayc"`
}

type SuitMappings struct {
	Masks    map[string]MaskMapping       `yaml:"masks"`
	SkipMask []string                     `yaml:"skip_mask"`
	Suits    map[string]map[string]string `yaml:"suits"`
}

type MaskMapping struct {
	Default       string            `yaml:"default"`
	TraitMappings map[string]string `yaml:"trait_mappings"`
	ChromaKey     string            `yaml:"chroma_key"`
}

type CoffeeMugMappings struct {
	Furs    map[string]string `yaml:"furs"`
	Liquids map[string]string `yaml:"liquids"`
	Steam   map[string]string `yaml:"steam"`
	Logos   map[string]string `yaml:"logos"`
}

type HoundTraitMappings struct {
	Faces  map[string]string `yaml:"faces"`
	Forms  map[string]string `yaml:"forms"`
	Heads  map[string]string `yaml:"heads"`
	Legs   map[string]string `yaml:"legs"`
	Mouths map[string]string `yaml:"mouths"`
	Noses  map[string]string `yaml:"noses"`
	Torsos map[string]string `yaml:"torsos"`
}

type BackgroundImagePlacementMappings struct {
	BAYCBackground  string `yaml:"bayc_background"`
	HoundBackground string `yaml:"hound_background"`
	MAYCBackground  string `yaml:"mayc_background"`
}

type BAYCBackgroundMappings struct {
	BAYCCornerMask          string              `yaml:"bayc_corner_mask"`
	BAYCCornerMaskChromaKey string              `yaml:"bayc_corner_mask_chroma_key"`
	BAYCBackgroundColorKeys map[string][]string `yaml:"bayc_background_color_keys"`
}

type ImageProcessorConfig struct {
	GMMappings                       map[string]string                `yaml:"gm_mappings"`
	NFDMerchMappings                 MerchMappings                    `yaml:"nfd_merch_mappings"`
	SuitMappings                     SuitMappings                     `yaml:"suit_mappings"`
	Hands                            map[string]string                `yaml:"-"`
	PledgeHands                      PledgeMappings                   `yaml:"pledge_hands"`
	ApeBagMappings                   map[string]string                `yaml:"ape_bag"`
	BAYCBackgroundMappings           BAYCBackgroundMappings           `yaml:"bayc_background_mappings"`
	MAYCBackgroundColorKeys          map[string]string                `yaml:"mayc_background_color_keys"`
	MAYCCoffeeMugMappings            CoffeeMugMappings                `yaml:"mayc_coffee_mug_mappings"`
	NFLJerseyMappings                map[string]map[string]string     `yaml:"nfl_jersey_mappings"`
	HoundTraitMappings               HoundTraitMappings               `yaml:"hound_trait_mappings"`
	SerumCityBackgroundImageMappings BackgroundImagePlacementMappings `yaml:"serum_city_background_image_mappings"`
	ApeCoinBackgroundImageMappings   BackgroundImagePlacementMappings `yaml:"apecoin_background_image_mappings"`
}

func LoadCfg(env, cfgFile string) (Config, error) {
	cfg := Config{}
	err := config.LoadWithEnvironmentPrefix(env, cfgFile, &cfg, validateConfig)
	return cfg, err
}

func validateConfig(cfg *Config) error {
	if len(cfg.HoundsMetadataEndpoint) == 0 {
		return fmt.Errorf("No Mutant Hounds Metadata Endpoint Defined")
	}

	if len(cfg.MAYCMetadataEndpoint) == 0 {
		return fmt.Errorf("No MAYC Metadata Endpoint Defined")
	}

	if len(cfg.BotToken) == 0 {
		return fmt.Errorf("No Discord Bot Token Defined")
	}

	if len(cfg.IPFSEndpoint) == 0 {
		return fmt.Errorf("No IPFS Endpoint Defined")
	}

	return nil
}
