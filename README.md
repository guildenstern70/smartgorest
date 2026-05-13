# Smart Go REST

[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![Go Report](https://goreportcard.com/badge/github.com/guildenstern70/golearn)](https://goreportcard.com/report/github.com/guildenstern70/golearn)

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

