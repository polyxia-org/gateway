package config

import (
	"log"

	"github.com/thomasgouveia/go-config"
)

const (
	DEFAULT_MORTY_CONTROLLER_ENDPOINT = "http://localhost:8083"
	DEFAULT_MORTY_REGISTRY_ENDPOINT   = "http://localhost:8081"
	DEFAULT_NLU_API_ENDPOINT          = "http://localhost:8082"
	DEFAULT_ADDRESS                   = "0.0.0.0"
	DEFAULT_PORT                      = 8080
)

type Config struct {
	Addr                    string `yaml:"addr"`
	Port                    int    `yaml:"port"`
	MortyControllerEndpoint string `yaml:"mortycontrollerendpoint"`
	MortyRegistryEndpoint   string `yaml:"mortyregistryendpoint"`
	NluApiEndpoint          string `yaml:"nluapiendpoint"`
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
			MortyRegistryEndpoint:   DEFAULT_MORTY_REGISTRY_ENDPOINT,
			MortyControllerEndpoint: DEFAULT_MORTY_CONTROLLER_ENDPOINT,
			NluApiEndpoint:          DEFAULT_NLU_API_ENDPOINT,
			Addr:                    DEFAULT_ADDRESS,
			Port:                    DEFAULT_PORT,
		},
	})
	if err != nil {
		log.Fatal("p", err)
	}

	return cl.Load()
}
