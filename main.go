package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	MORTY_API_ENDPOINT := getEnv("MORTY_API_ENDPOINT", "http://localhost:8081")
	// NLU_API_ENDPOINT := getEnv("NLU_API_ENDPOINT", "http://localhost:8081")

	router := gin.Default()

	router.POST("/skills", func(c *gin.Context) {
		/* Get name param */
		name := c.PostForm("name")
		log.Printf("name: %s", name)
		// lowercase name to match morty registry compliance
		name = strings.ToLower(name)

		mortyFunctionRegistryResp, err := handleArchive(MORTY_API_ENDPOINT, c, name)
		if err != nil {
			log.Printf("handleArchive")
			log.Fatal(err)
		}

		if mortyFunctionRegistryResp.StatusCode != http.StatusOK {
			log.Printf(mortyFunctionRegistryResp.Status)
			log.Fatal(err)
		}
		/*
			TODO: handle intents_json
			nluResp, err := handleIntentsJSON(NLU_API_ENDPOINT, c, name)
			if err != nil {
				log.Printf("handleIntentsJSON")
				log.Fatal(err)
			}

			if nluResp.StatusCode != http.StatusCreated {
				log.Printf(nluResp.Status)
				log.Fatal(err)
			} */

		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Skill successfully created!"})
	})

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

// getEnv returns the value of the environment variable named by the key.
// If the variable is not present, it returns the `fallback` value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func handleIntentsJSON(NLU_API_ENDPOINT string, c *gin.Context, name string) (*http.Response, error) {
	// Get intentsJsonFile from intents from request (formfile: intentsJsonFile)
	intentsJsonFile, _, err := c.Request.FormFile("intents_json")
	if err != nil {
		log.Printf("Read intents_json from request (formfile: intents_json)")
		log.Fatal(err)
	}
	defer intentsJsonFile.Close()
	intentsJson := new(bytes.Buffer)
	intentsJson.ReadFrom(intentsJsonFile)
	// POST the json
	/*
		curl -X POST \
		http://localhost:8082/v1/intents \
		-H 'Content-Type: application/json' \
		-d '{
			"name": "lighton",
			"intents_json": {
				"intent": "LightOn",
				"utterances": [
				"lumière",
				"allume la lumière",
				"aziz lumière",
				"allume la lampe"
			}
		}'
	*/
	// Build json body
	jsonBody := new(bytes.Buffer)
	jsonBody.WriteString(`{"name":"`)
	jsonBody.WriteString(name)
	jsonBody.WriteString(`","intents_json":`)
	jsonBody.Write(intentsJson.Bytes())
	jsonBody.WriteString(`}`)

	// Send intentsJson to NLU
	req, err := http.NewRequest("POST", NLU_API_ENDPOINT+"/v1/intents", jsonBody)
	if err != nil {
		log.Printf("NewRequest")
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send intentsJson to NLU
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do")
		log.Fatal(err)
	}

	return resp, err
}

func handleArchive(MORTY_API_ENDPOINT string, c *gin.Context, name string) (*http.Response, error) {
	// Read function rootfs file from request (formfile: function_archive)
	file, _, err := c.Request.FormFile("function_archive")
	if err != nil {
		log.Printf("Read function rootfs file from request (formfile: function_archive)")
		log.Fatal(err)
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
		log.Printf("CreateFormFile")
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Printf("Copy")
		log.Fatal(err)
	}
	w.Close()

	// Send function_archive to Morty Function Registry
	req, err := http.NewRequest("POST", MORTY_API_ENDPOINT+"/v1/functions/build", &b)
	if err != nil {
		log.Printf("NewRequest")
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send function_archive to Morty Function Registry
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return resp, err
}
