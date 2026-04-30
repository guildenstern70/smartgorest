# SMART-GO-REST PROJECT

[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![Go Report](https://goreportcard.com/badge/github.com/guildenstern70/golearn)](https://goreportcard.com/report/github.com/guildenstern70/golearn)

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

## Ecosystem

•	Framework (Echo): Highly recommended for REST. It is more "feature-complete" than Gin and feels more 
    idiomatic for building professional APIs.
•	OpenAPI (Swag): You write comments over your functions to generate the spec (Code-first).
•	ORM (Ent.): Created by Facebook, Ent is a "graph-based" ORM. It generates a type-safe API for your 
    database schema. It is significantly more robust and less "magic" (using reflection) than GORM, 
    making it feel more like a modern tool.
Why it fits:
•	Speed of Development: You can get a Go API running in minutes.
•	Explicit over Implicit: If you found Quarkus too "magical," you'll love Go's clarity.
•	Ent's Schema: Defining models in Ent feels very similar to defining Pydantic models or Java Entities.


### Add modules to the project

    go get -u [fully qualified module name]

### Ent.

Add a new entity

    go run -mod=mod entgo.io/ent/cmd/ent new User

Clear the database and regenerate the code

    go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema --feature sql/lock

Notice that the changes are applied to the database when you run the app, because
it has a mechanism to auto-create the schema on startup.

