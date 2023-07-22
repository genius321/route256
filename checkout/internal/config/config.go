package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	GrpcPort    int    `yaml:"grpcPort"`
	HttpPort    int    `yaml:"httpPort"`
	ServiceName string `yaml:"serviceName"`
	Services    struct {
		Loms           string `yaml:"loms"`
		ProductService string `yaml:"productService"`
		Pprof          int    `yaml:"pprof"`
	} `yaml:"services"`
	RateLimit int `yaml:"rateLimit"`
}

var AppConfig = Config{}

func Init() error {
	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return nil
}
