# gateway

The goal of this project is to provide a simple API gateway for **creating skills**.

The gateway:
- Creates a new functions with the Morty Function Registry
- Creates a new skill with the NLU API

## Prerequisites

* [Golang](https://go.dev/doc/install) (`>=1.20`)
* Docker and Docker Compose (for running the Morty Function Registry locally)

## Getting Started

Flow:

![Flow](./docs/flow.svg)

### Run the API gateway

1. Set the environment variables for the NLU and Morty APIs:

```bash
# Example values
export POLYXIA_GATEWAY_MORTYREGISTRYENDPOINT="http://localhost:8081"
export POLYXIA_GATEWAY_MORTYCONTROLLERENDPOINT="http://localhost:8081"
export POLYXIA_GATEWAY_NLUAPIENDPOINT="http://localhost:8082"
export POLYXIA_GATEWAY_ADDR="localhost"
export POLYXIA_GATEWAY_PORT="8080"
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

Follow the instructions here: https://github.com/polyxia-org/nlu/

Example instructions:
```bash
cd ..
git clone https://github.com/polyxia-org/nlu
cd nlu
git pull -r
docker run -d --rm --pull always -p 8000:8000 surrealdb/surrealdb:latest start --pass root
poetry run python nlu/app.py
```

### Create a new skill

When creating new skills. You have to know that the intent name should be chosen wisely. The intent name will be used to make the right choice of function to use in the NLU.
The intent name should be explicit and use underscores to separate words.
example of good intent name: `get_weather`, `get_quote`, `say_hello`, `order_food`, `today_date`

Also, you can use query parameters those are the slots in your skills. For example, `get_weather` can use a query parameter `place_name` to get the weather of a specific place.
There is a list of query parameters that you can use in your skills. But you cannot use another query parameter that is not in the list.
Available slots type list is [here](./docs/slot_types.txt).

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

For more information, see the [OpenAPI spec](./openapi.yml).

### Run the skill

1. Run the skill using the Morty CLI:

```bash
curl -X POST \
  http://localhost:8080/v1/nlu \
  -H 'Content-Type: application/json' \
  -d '{"input_text": "quelle est la météo ?"}'
```
