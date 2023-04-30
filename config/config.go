package config

type Config struct {
	MetadataEndpoint     string         `yaml:"metadata_endpoint"`
	ImageProcessorConfig ImageProcessor `yaml:"image_processor_mappings"`
}

type ImageProcessor struct {
	GMMappings map[string]string `yaml:"gm_mappings"`
}
