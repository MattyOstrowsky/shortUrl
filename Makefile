# Default variables
SERVER_PORT		?= 2345
MIN_LENGTH		?= 6
MAX_LENGTH		?= 32
DEFAULT_LENGTH	?= 8
DB_NAME			?= urldb
DB_USER			?= root
DB_PASS			?= example
DB_PORT			?= 27017
COLLECTION_NAME	?= urls
COMPOSE_ARGS 	:= MIN_LENGTH=${MIN_LENGTH} MAX_LENGTH=${MAX_LENGTH} DEFAULT_LENGTH=${DEFAULT_LENGTH} SERVER_PORT=${SERVER_PORT} DB_NAME=${DB_NAME} DB_USER=${DB_USER} DB_PASS=${DB_PASS} DB_PORT=${DB_PORT} COLLECTION_NAME=${COLLECTION_NAME}
COMPOSE 		:= ${COMPOSE_ARGS} docker compose

# Targets for running MongoDB
run-mongodb:
	@echo "Running MongoDB ..."
	@${COMPOSE} up -d mongodb
	@echo "MongoDB is running."

# Targets for running URL API
run-urlapi:
	@echo "Running URL API ..."
	@${COMPOSE} up -d urlapi
	@echo "URL API is running."

# Targets for building URL CLI
build-urlcli:
	@echo "Building URL CLI ..."
	@${COMPOSE} build urlcli
	@echo "URL CLI is built."

# Targets for running URL CLI
run-urlcli:
	@${COMPOSE} run urlcli ${CMD}

# Targets for running URL API locally
local-run-urlapi:
	cd src/urlapi && \
	go mod download && \
	go build -o ../../bin/urlapi ./main.go && \
	cp ./template.html ../../bin/template.html && \
	../../bin/urlapi

# Targets for running URL CLI locally
local-urlcli:
	cd src/urlcli && \
	go mod download && \
	go build -o ../../bin/urlcli ./main.go && \
	../../bin/urlcli 

.PHONY: run-mongodb run-urlapi build-urlcli run-urlcli local-run-urlapi local-urlcli
