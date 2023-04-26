package gateway

import (
	"encoding/json"
	"net/http"
)

type healthcheck struct {
	Status string `json:"status"`
}

const (
	up           = "UP"
	outOfService = "OUT_OF_SERVICE"
)

func (s *Server) HealthcheckHandler(w http.ResponseWriter, _ *http.Request) {

	// check is morty is up
	_, err := http.Get(s.cfg.MortyRegistryEndpoint + "/healthz")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&healthcheck{Status: outOfService})
		return
	}

	// check is rick is up
	_, err = http.Get(s.cfg.NluApiEndpoint + "/docs")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&healthcheck{Status: outOfService})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&healthcheck{Status: up})
}
