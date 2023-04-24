package gateway

import (
	"encoding/json"
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
	responseBody := SendToNLU(s.cfg.NluApiEndpoint, r.Body)

	// Create the response object
	nluResponse := NLUResponse{Body: responseBody}

	// Marshal the response object
	responseBytes, err := json.Marshal(nluResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling NLUResponse object."))
		return
	}

	// Send the response to the device
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func SendToNLU(NLU_API_ENDPOINT string, inputBody io.Reader) string {
	if inputBody != nil {
		bodyBytes, err := ioutil.ReadAll(inputBody)
		if err != nil {
			return err.Error()
		}
		log.Debugf("bodyBytes: %s", string(bodyBytes))
	}

	// Send intentsJson to NLU
	req, err := http.NewRequest("POST", NLU_API_ENDPOINT+NLU_PATH, inputBody)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")

	// Send intentsJson to NLU
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
