package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      string `yaml:"port"`
	TimeStamp string `yaml:"timeStamp"`
}

const fileName = "config.yml"

func NewConfig() (Config, error) {
	config := Config{Port: "8000"}

	b, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to open config file: ", fileName)
		return config, err
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Println("Failed to Unmarshal config file: ", fileName)
		return config, err
	}

	return config, nil
}
