# gateway

The goal of this project is to provide a simple API gateway for **creating skills**.

By recieving a POST request like this:

```bash
curl -X POST \
  http://localhost:8080/v1/skills \
  -H 'Content-Type: multipart/form-data' \
  -F 'name=lighton' \
  -F 'intents_json=@./test_data/intent.json' \
  -F 'function_archive=@./test_data/lightOn.zip'
```

The gateway will:
- Create a new functions with the Morty Function Registry
- Create a new skill with the NLU API

**Today, the gate do not "Create a new skill with the NLU API"** because the NLU API is not ready yet. 

## Prerequisites

* [Golang](https://go.dev/doc/install) (`>=1.20`)

## Getting Started

Flow without NLU:

![Flow](./docs/flow.svg)

### Run the API gateway

1. Set the environment variables for the NLU and Morty APIs:

```bash
# Example values
export MORTY_API_ENDPOINT="http://localhost:8081/v1/functions/build"
export NLU_API_ENDPOINT="http://localhost:8082/v1/skills"
```

2. Run the API gateway with the following command:
```bash
go run main.go
```

3. Run Morty Function Registry:
```bash
cd ..
git clone https://github.com/polyxia-org/morty-registry
cd morty-registry
git pull -r
docker compose up -d
export AWS_ACCESS_KEY_ID=mortymorty
export AWS_SECRET_ACCESS_KEY=mortymorty

# create 'functions' bucket
aws --endpoint-url=http://localhost:9000 s3 mb s3://functions

# will ask password for running sudo
make start
```

4. Run the NLU API
```bash
# TODO when the NLU API is ready
```

### Use the API gateway

1. Create a new skill using the Morty CLI:
    
```bash
cd ./test_data
# Install the CLI by following the instructions here
# download from https://github.com/polyxia-org/morty-cli/releases/download/v1.0.0/morty-v1.0.0-linux-amd64.tar.gz
tar -xvf morty.tar.gz
sudo mv morty /usr/local/bin

# Create a new skill
morty function init lightOn
zip -r lightOn.zip lightOn
```

2. Write an `intent.json` to `./test_data`

The json should look like this:
```json
{
    "utterances": [
      "météo",
      "donne moi la météo",
      "quel temps fait-il aujourd'hui à [ Paris | Montpellier ] ?",
      "est ce qu'il pleut ?",
      "est-ce qu'il fait beau ?",
      "Y a t'il du soleil ?",
      "Quelle est la météo actuelle ?"
    ],
    "slots": [
      {
        "type": "place_name"
      }
    ]
  }
```

3. Send a POST request with 2 files to the API gateway:

```bash
curl -X POST \
  http://localhost:8080/v1/skills \
  -H 'Content-Type: multipart/form-data' \
  -F 'name=lighton' \
  -F 'intents_json=@./test_data/intent.json' \
  -F 'function_archive=@./test_data/lightOn.zip'
```