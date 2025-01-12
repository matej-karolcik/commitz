package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Model       string  `yaml:"model" mapstructure:"model" default:"llama3.2"`
	Temperature float64 `yaml:"temperature" mapstructure:"temperature" default:"0.5"`
}

func Load(path ...string) (*Config, error) {
	viper.SetConfigType("yaml")

	if len(path) == 0 {
		viper.SetConfigName(".commitz")
		workDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		viper.AddConfigPath(workDir)
	} else {
		viper.SetConfigFile(path[0])
	}

	_ = viper.ReadInConfig()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	defaults.SetDefaults(&config)

	return &config, nil
}

func Dump(path ...string) error {
	config, err := Load(path...)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println(string(yamlData))

	var filePath string
	if len(path) > 0 {
		filePath = path[0]
	} else {
		filePath, err = getConfigFile()
		if err != nil {
			return fmt.Errorf("failed to get config file: %w", err)
		}
	}

	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func getConfigFile() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	return filepath.Join(workDir, ".commitz.yaml"), nil
}
