package config

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type ValidatorFn[T any] func(configObj *T) error

func Load[T any](path string, dest *T, ValidatorFns ...ValidatorFn[T]) error {
	return LoadWithEnvironmentPrefix("", path, dest, ValidatorFns...)
}

func LoadWithEnvironmentPrefix[T any](envPrefix, path string, dest *T, ValidatorFns ...ValidatorFn[T]) error {
	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Error Reading Config: %w", err)
	}

	cfg := make(map[string]interface{})
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("Error Parsing Config: %w", err)
	}

	//Viper doesn't respect case sensitivty so
	//as a workaround unmarshal the yaml config into a map
	//and then use the mapstructure decoder library
	//to unamrshal the map into the config struct
	//using the appropriate yaml tags
	decoderCfg := &mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "yaml",
	}

	cfgDecoder, err := mapstructure.NewDecoder(decoderCfg)
	if err != nil {
		return fmt.Errorf("Error Instantating Config Decoder: %w", err)
	}

	if err := cfgDecoder.Decode(cfg); err != nil {
		return fmt.Errorf("Error Decoding Config: %w", err)
	}

	for _, validate := range ValidatorFns {
		if err := validate(dest); err != nil {
			return fmt.Errorf("Error valdiating config: %w", err)
		}
	}

	return nil
}
