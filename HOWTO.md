# Url Shortener

This project consists of 3 services/applications:
- urlcli: A CLI for managing URLs.
- urlapi: An API for URL shortening, used for redirection to short URLs.
- mongodb: The database used by the above services.

##  Requirements:
- Docker
- Docker Compose v2
- Make
- Go

## Setup
The makefile and docker/docker-compose-v2 are used for managing applications. Application settings are configured using environment variables. It is important to set the following environment variables before running the application for it to work correctly.

wymagania:
docker 
docker compose v2
make
go


```
export SERVER_PORT=2345     # Urlapi port.
export MIN_LENGTH=6         # Minimum length for generated short URLs.
export MAX_LENGTH=32        # Maximum length for generated short URLs.
export DEFAULT_LENGTH=8     # Default length for generated short URLs.

# mongodb variables
export DB_NAME=urldb         # Name of the MongoDB database.
export DB_USER=root          # MongoDB username. 
export DB_PASS=example       # MongoDB password. 
export DB_PORT=27017         # MongoDB port.
export COLLECTION_NAME=urls  # Name of the collection in the MongoDB database.
```

You can run the applications either in a Docker environment or locally.

#### MongoDB Setup

For proper functionality, MongoDB database is required.
To set up MongoDB, run:
```
make run-mongodb
```

### Local
***Note: Before running the applications, ensure that the necessary environment variables are set for proper configuration and operation.***

#### Urlapi Setup

To run the API server, execute:
```
make local-run-urlapi
```

Check if the server is running:
```
curl http://localhost:${SERVER_PORT}
```

#### Urlcli Setup

To set up urlcli, run:
```
make local-urlcli
```

You can either export the PATH:
```
export PATH="$(pwd)/bin:$PATH"
```

And use the tool:
```
urlcli -h
```

Or execute the binary directly:
```
./bin/urlcli -h
```

### Docker
***Note: Before running the applications, ensure that the necessary environment variables are set for proper configuration and operation.***

#### Urlapi Setup

To run the API server, execute:
```
make run-urlapi
```

Check if the server is running:
```
curl http://localhost:${SERVER_PORT}
```

#### Urlcli Setup

To set up urlcli, run:
```
make build-urlcli
```

Now, you can use the image:
```
make run-urlcli CMD="add http://example.com"
```

where CMD is the command to execute.

Alternatively, you can run the container directly:
```
docker compose run urlcli [CMD]
# Example:
# docker compose run urlcli add http://example.com
# docker compose run urlcli list
```

## Usage

### urlapi

Once the API is up and running, it will be available at http://localhost:[SERVER_PORT].
Shortened URLs can be accessed at http://localhost:[SHORTENED_URL], for example, http://localhost:XYZ123.
You will be redirected to the correct website immediately upon accessing the link.
If you don't want to be redirected immediately, you can use the ! sign before SHORTENED_URL, e.g., http://localhost:!XYZ123. In this case, a redirect link to the website will be returned.


### urlcli

To see available commands and flags, use:
```
urlcli -h
```

You can set the maximum length parameter of the shortened link either by default using an environment variable or by overriding it with the -ml flag:
```
urlcli add http://example.com --lm=10
```

Some examples:
```
urlcli add http://example.com --lm=10 -v
urlcli update http://example.com --lm=12 -v
urlcli list
urlcli delete http://example.com
```