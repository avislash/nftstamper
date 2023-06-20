package config

import (
	"fmt"

	"github.com/avislash/nftstamper/lib/config"
	"github.com/avislash/nftstamper/lib/filter"
)

type Config struct {
	MetadataEndpoint     string               `yaml:"metadata_endpoint"`
	BotToken             string               `yaml:"discord_bot_token"`
	IPFSEndpoint         string               `yaml:"ipfs_endpoint"`
	LogLevel             string               `yaml:"log_level"`
	ImageProcessorConfig ImageProcessorConfig `yaml:"image_processor_mappings"`
}

type ImageProcessorConfig struct {
	GMMappings        map[string]map[string]string `yaml:"gm_mappings"`
	BaseGMSmoke       string                       `yaml:"base_gm_smoke"`
	BaseGMSmokePath   string                       `yaml:"base_gm_smoke_path"`
	AzulGMSmoke       string                       `yaml:"azul_gm_smoke"`
	BaseGMSmokeBorder string                       `yaml:"base_gm_smoke_border"`
	AzulGMSmokeBorder string                       `yaml:"azul_gm_smoke_border"`
	HBDMappings       map[string]string            `yaml:"hbd_mappings"`
	Filters           FilterCfgs                   `yaml:"filters"`
}

type FilterCfgs struct {
	Opacity map[string]filter.Filter[string, float64] `yaml:"opacity"`
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
