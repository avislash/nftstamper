package config

import (
	"fmt"

	"github.com/avislash/nftstamper/lib/config"
)

type Config struct {
	HoundsMetadataEndpoint string               `yaml:"hounds_metadata_endpoint"`
	MAYCMetadataEndpoint   string               `yaml:"mayc_metadata_endpoint"`
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

type ImageProcessorConfig struct {
	GMMappings       map[string]string `yaml:"gm_mappings"`
	NFDMerchMappings MerchMappings     `yaml:"nfd_merch_mappings"`
	Suit             string            `yaml:"suit"`
	Hands            map[string]string `yaml:"-"`
	PledgeHands      PledgeMappings    `yaml:"pledge_hands"`
	ApeBagMappings   map[string]string `yaml:"ape_bag"`
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
