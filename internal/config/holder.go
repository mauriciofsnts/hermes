package config

import (
	"os"

	"github.com/ghodss/yaml"
)

var (
	Hermes *Config
)

func LoadConfig() error {
	f, err := os.ReadFile("./config.yml")

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, &Hermes)

	if err != nil {
		return err
	}

	return nil
}
