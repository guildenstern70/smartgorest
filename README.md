# Smart Go REST

[![License: ISC](https://img.shields.io/badge/License-ISC-blue.svg)](https://opensource.org/licenses/ISC)
[![Go Report Card](https://goreportcard.com/badge/github.com/guildenstern70/smartgorest)](https://goreportcard.com/report/github.com/guildenstern70/smartgorest)

<img src="/public/screenshot.png" width="800"  alt="Smart Go REST screenshot"/>

## Setup

You should have a single workspace repository for all of your GO
code, and a GOPATH pointing to that.

It is advisable to have a /src directory inside your workspace.

This code must be installed in a directory inside

    $GOPATH/src

with the following path

    github.com/guildenstern70/smartgorest

so that you find the LICENSE file here:

    $GOPATH/src/github.com/guildenstern70/smartgorest/LICENSE

## Build

First, update the dependencies by running

    go mod tidy

You build the executable file by running

    go install .

## Run

    go run .

To run the executable file:

    $GOPATH/bin/smartgorest.exe (Windows)
    $GOPATH/bin/smartgorest (Mac + Linux)

## Test
Run the test suite of this project by running

    go test github.com/guildenstern70/smartgorest/internal

## Docker / Podman

This project includes a Dockerfile and docker-compose configuration to containerize the application with PostgreSQL.

### Building with Docker

To build the Docker image:

    docker build -t smartgorest:latest .

### Building with Podman

To build the image with Podman:

    podman build -t smartgorest:latest .

### Running with Docker Compose

The easiest way to run the application with PostgreSQL is using docker-compose:

    docker-compose up

This will:
- Build the application image
- Start a PostgreSQL 16 database
- Start the SmartGoRest application
- Expose the API on http://localhost:1323

To run in the background:

    docker-compose up -d

To stop the containers:

    docker-compose down

To view logs:

    docker-compose logs -f app

### Running with Podman Compose

If you prefer Podman, use podman-compose (install it first if needed):

    pip install podman-compose

Then run:

    podman-compose up

All other commands work the same as Docker Compose (substitute `podman-compose` for `docker-compose`).

### Running Container Standalone

To run the application container standalone with Docker:

    docker run -p 1323:1323 \
      -e DB_HOST=<postgres_host> \
      -e DB_PORT=<postgres_port> \
      -e DB_USER=<db_user> \
      -e DB_PASSWORD=<db_password> \
      -e DB_NAME=<db_name> \
      smartgorest:latest

With Podman:

    podman run -p 1323:1323 \
      -e DB_HOST=<postgres_host> \
      -e DB_PORT=<postgres_port> \
      -e DB_USER=<db_user> \
      -e DB_PASSWORD=<db_password> \
      -e DB_NAME=<db_name> \
      smartgorest:latest

Make sure you have a PostgreSQL database running and accessible at the specified host and port.

### Environment Variables

Configure the application using environment variables:

- `DB_HOST`: PostgreSQL host (default: `postgres` in docker-compose)
- `DB_PORT`: PostgreSQL port (default: `5432`)
- `DB_USER`: PostgreSQL user (default: `smartgorest`)
- `DB_PASSWORD`: PostgreSQL password (default: `smartgorest`)
- `DB_NAME`: PostgreSQL database name (default: `smartgorest`)
- `DATABASE_URL`: Full connection string (auto-constructed from above variables if not provided)

### Environment File

Create a `.env` file in the project root to override default environment variables:

    DB_USER=myuser
    DB_PASSWORD=mypassword
    DB_NAME=mydatabase

Copy from `.env.example`:

    cp .env.example .env

## Ecosystem

Read here about the ecosystem of this project:
https://medium.com/@alessiosaltarin/beyond-the-jvm-building-truly-native-microservices-with-go-83f383d13455


### Add modules to the project

    go get -u [fully qualified module name]

### Ent.

Add a new entity

    go run -mod=mod entgo.io/ent/cmd/ent new Person

Regenerate the code

    go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema --feature sql/lock

Notice that the changes are applied to the database when you run the app, because
it has a mechanism to auto-create the schema on startup.

## Swag

Refresh the API documentation

    swag init -g main.go

