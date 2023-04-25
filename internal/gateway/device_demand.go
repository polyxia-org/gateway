package gateway

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	NLU_PATH = "/v1/nlu"
)

// Payload is the input data structure for the request
type Payload struct {
	Payload string `json:"input_text"`
}

// NLUResponse is the output data structure for the response
type NLUResponse struct {
	Body string `json:"response"`
}

func (s *Server) DeviceDemandHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed."))
		return
	}

	// Send the payload to NLU service
	responseBody, err := getRespFromNLU(s.cfg.NluApiEndpoint, r.Body)
	if err != nil {
		log.Errorf("Error getting response from NLU: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting response from NLU."))
		return
	}

	// Send the response to the device
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func getRespFromNLU(nluApiEndpoint string, inputBody io.Reader) ([]byte, error) {
	if inputBody != nil {
		// only for debugging
		bodyBytes, err := ioutil.ReadAll(inputBody)
		if err != nil {
			return nil, err
		}
		log.Debugf("bodyBytes: %s", string(bodyBytes))
	}

	// Send json to NLU
	req, err := http.NewRequest("POST", nluApiEndpoint+NLU_PATH, inputBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBodyBytes, nil
}
