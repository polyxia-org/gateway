package main

import (
	"github.com/polyxia-org/gateway/internal/gateway"
	"github.com/polyxia-org/gateway/pkg/helpers"
	log "github.com/sirupsen/logrus"
)

const (
	LOG_LEVEL = "GATEWAY_LOG"
)

func main() {
	// Init logger for the app
	level := log.InfoLevel
	envLevel := helpers.GetEnv(LOG_LEVEL, "INFO")
	if envLevel != "" {
		lvl, err := log.ParseLevel(envLevel)
		if err != nil {
			log.Fatalf("failed to parse log level from environment: %v", err)
		}
		level = lvl
	}
	log.SetLevel(level)

	// Run the registry HTTP server
	gw, err := gateway.NewServer()
	if err != nil {
		log.Fatalf("failed to initialize the gateway: %v", err)
	}

	gw.Serve()
}
