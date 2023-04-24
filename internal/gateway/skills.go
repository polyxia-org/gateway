package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	MORTY_FUNCTIONS_BUILD_ENDPOINT = "/functions/build"
	NLU_SKILLS_ENDPOINT            = "/skills"
)

var (
	ErrRequiredFunctionName    = errors.New("the function name is required")
	ErrRequiredFunctionRuntime = errors.New("the function runtime is required")
	ErrInvalidFunctionArchive  = errors.New("the function code archive must be a valid zip file")
)

func (s *Server) SkillsHandler(w http.ResponseWriter, r *http.Request) {

	/* Get name param */
	name := r.PostFormValue("name")
	// lowercase name to match morty registry compliance
	name = strings.ToLower(name)
	log.Debugf("name: %s", name)

	nluResp, err := handleIntentsJSON(s.cfg.NluApiEndpoint, r, name)
	if err != nil {
		log.Debugf("handleIntentsJSON", nluResp.Status)
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
		return
	}

	if nluResp.StatusCode != http.StatusCreated {
		log.Errorf(nluResp.Status)
		bodyBuf := new(bytes.Buffer)
		bodyBuf.ReadFrom(nluResp.Body)
		log.Debugf(bodyBuf.String())
		s.JSONResponse(w, nluResp.StatusCode, bodyBuf.String())
		return
	}

	mortyFunctionRegistryResp, err := handleArchive(s.cfg.MortyApiEndpoint, r, name)
	if err != nil {
		log.Error(err)
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
	}

	if mortyFunctionRegistryResp.StatusCode != http.StatusOK {
		log.Warnf("Morty registry creation failed!")
		// undo NLU skill creation with DELETE /v1/skills/:name
		req, _ := http.NewRequest("DELETE", s.cfg.NluApiEndpoint+NLU_SKILLS_ENDPOINT+"/"+name, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		_, err := client.Do(req)
		if err != nil {
			log.Debugf("DELETE", s.cfg.NluApiEndpoint+NLU_SKILLS_ENDPOINT+"/"+name)
		}
		log.Debugf(mortyFunctionRegistryResp.Status)
		bodyBuf := new(bytes.Buffer)
		bodyBuf.ReadFrom(mortyFunctionRegistryResp.Body)
		log.Debugf(bodyBuf.String())
		s.JSONResponse(w, mortyFunctionRegistryResp.StatusCode, bodyBuf.String())
		return
	}

	s.JSONResponse(w, http.StatusOK, APISucess{
		StatusCode: http.StatusOK,
		Message:    "Skill created successfully",
	})
}

func makeAPIError(status int, err error) *APIError {
	return &APIError{
		StatusCode: status,
		Message:    err.Error(),
	}
}

func handleIntentsJSON(NLU_API_ENDPOINT string, r *http.Request, name string) (*http.Response, error) {
	// Get intentsJsonFile from intents from request (formfile: intentsJsonFile)
	intentsJsonFile, _, err := r.FormFile("intents_json")
	if err != nil {
		return nil, err
	}
	defer intentsJsonFile.Close()
	intentsJson := new(bytes.Buffer)
	intentsJson.ReadFrom(intentsJsonFile)
	// Map Byte Buffer to JSON
	var intents map[string]interface{}
	// Unmarshal the JSON data into the intents map
	json.Unmarshal(intentsJson.Bytes(), &intents)
	// Add intent name to ensure the same name is used in NLU and Morty
	intents["intent"] = name
	// Convert the modified Go value back to JSON format
	intentsJson.Reset()
	if err := json.NewEncoder(intentsJson).Encode(intents); err != nil {
		return nil, err
	}

	// Send intentsJson to NLU
	req, err := http.NewRequest("POST", NLU_API_ENDPOINT+NLU_SKILLS_ENDPOINT, intentsJson)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send intentsJson to NLU
	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

func handleArchive(MORTY_API_ENDPOINT string, r *http.Request, name string) (*http.Response, error) {
	// Read function rootfs file from request (formfile: function_archive)
	file, _, err := r.FormFile("function_archive")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	/*
		curl -X POST \
		http://localhost:8081/v1/functions/build \
		-H 'Content-Type: multipart/form-data' \
		-F 'name=lighton' \
		-F 'runtime=node-19' \
		-F 'archive=@./test_data/lightOn.zip'
	*/
	// The function_archive file should be sent to Morty Function Registry
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add form fields
	w.WriteField("name", name)
	w.WriteField("runtime", "node-19")

	fw, err := w.CreateFormFile("archive", name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	w.Close()

	// Send function_archive to Morty Function Registry
	req, err := http.NewRequest("POST", MORTY_API_ENDPOINT+MORTY_FUNCTIONS_BUILD_ENDPOINT, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send function_archive to Morty Function Registry
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyString := string(bodyBytes)
		log.Debugf(bodyString)
	}

	return resp, nil
}
