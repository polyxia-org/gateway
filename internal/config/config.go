package config

import (
	"log"

	"github.com/thomasgouveia/go-config"
)

const (
	DEFAULT_MORTY_API_ENDPOINT = "http://localhost:8081"
	DEFAULT_NLU_API_ENDPOINT   = "http://localhost:8082"
	DEFAULT_ADDRESS            = "0.0.0.0"
	DEFAULT_PORT               = 8080
)

type Config struct {
	Addr             string `yaml:"addr"`
	Port             int    `yaml:"port"`
	MortyApiEndpoint string `yaml:"morty_api_endpoint"`
	NluApiEndpoint   string `yaml:"nlu_api_endpoint"`
}

func Load() (*Config, error) {
	cl, err := config.NewLoader(&config.Options[Config]{
		// You can use config.JSON also if you prefer
		Format: config.YAML,

		// Configuration file
		FileName:      "polyxia-gateway",
		FileLocations: []string{"/etc/polyxia-gateway", "."}, // Will search for a "my-application.yaml" file into the directories

		// Enable automatic environment variables lookup
		EnvEnabled: true,
		EnvPrefix:  "POLYXIA_GATEWAY",

		Default: &Config{
			MortyApiEndpoint: DEFAULT_MORTY_API_ENDPOINT,
			NluApiEndpoint:   DEFAULT_NLU_API_ENDPOINT,
			Addr:             DEFAULT_ADDRESS,
			Port:             DEFAULT_PORT,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return cl.Load()
}
