package config

import (
	"fmt"

	"github.com/avislash/nftstamper/lib/config"
)

type Config struct {
	MetadataEndpoint     string               `yaml:"metadata_endpoint"`
	BotToken             string               `yaml:"discord_bot_token"`
	IPFSEndpoint         string               `yaml:"ipfs_endpoint"`
	ImageProcessorConfig ImageProcessorConfig `yaml:"image_processor_mappings"`
}

type ImageProcessorConfig struct {
	GMMappings map[string]string `yaml:"gm_mappings"`
}

func LoadCfg(env, cfgFile string) (Config, error) {
	cfg := Config{}
	err := config.LoadWithEnvironmentPrefix(env, cfgFile, &cfg, validateConfig)
	return cfg, err
}

func validateConfig(cfg *Config) error {
	if len(cfg.MetadataEndpoint) == 0 {
		return fmt.Errorf("No Metadata Endpoint Defined")
	}

	if len(cfg.BotToken) == 0 {
		return fmt.Errorf("No Discord Bot Token Defined")
	}

	if len(cfg.IPFSEndpoint) == 0 {
		return fmt.Errorf("No IPFS Endpoint Defined")
	}

	return nil
}
