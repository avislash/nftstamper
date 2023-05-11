package config

type Config struct {
	MetadataEndpoint     string               `yaml:"metadata_endpoint"`
	IPFSEndpoint         string               `yaml:"ipfs_endpoint"`
	ImageProcessorConfig ImageProcessorConfig `yaml:"image_processor_mappings"`
}

type ImageProcessorConfig struct {
	GMMappings map[string]string `yaml:"gm_mappings"`
}
