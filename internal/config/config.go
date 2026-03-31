package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CapabilityOption struct {
	Provider string `yaml:"provider,omitempty"`
	Path     string `yaml:"path,omitempty"`
}

type Config struct {
	Gateway           string                      `yaml:"gateway"`
	Port              int                         `yaml:"port"`
	Token             string                      `yaml:"token"`
	TLS               bool                        `yaml:"tls"`
	Discovery         string                      `yaml:"discovery"`
	Capabilities      map[string]bool             `yaml:"capabilities"`
	CapabilityOptions map[string]CapabilityOption `yaml:"capabilityOptions,omitempty"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Capabilities == nil {
		cfg.Capabilities = Default().Capabilities
	}
	if cfg.CapabilityOptions == nil {
		cfg.CapabilityOptions = map[string]CapabilityOption{}
	}
	return &cfg, nil
}

func Default() *Config {
	return &Config{
		Port:      18789,
		Discovery: "auto",
		Capabilities: map[string]bool{
			"camera":        true,
			"location":      true,
			"photos":        true,
			"screen":        true,
			"motion":        false,
			"notifications": true,
			"sms":           false,
			"calendar":      false,
		},
		CapabilityOptions: map[string]CapabilityOption{},
	}
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
