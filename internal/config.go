package internal

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config represents the structure of global configuration
type Config struct {
	UserID string `yaml:"user_id"`
	Token  string `yaml:"token"`
}

// CreateDefaultConfig creates an empty configuration
func CreateDefaultConfig() *Config {
	return &Config{}
}

// LoadConfig loads the configuration from a YAML file or creates the empty one
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		defaultConf := CreateDefaultConfig()
		_ = SaveConfig(defaultConf, filename)
		return defaultConf, nil
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to a YAML file
func SaveConfig(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
