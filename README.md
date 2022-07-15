# speakeasy-example-rest-service-go

## Getting Started

### Prerequisites

* Go 1.18 (though should be backwards compatible with earlier versions)

### Running locally

1. From root of the repo
2. Run `go mod download` to install dependencies
3. Run `docker-compose up -d postgres` to run the postgres dependency
4. Run `go run cmd/server/main.go` will start the server on port 8080

### Running via docker

1. From root of the repo
2. Run `docker-compose up` will start the dependencies and server on port 8080

### Postman

I have provided Postman collections for testing out the REST endpoints exposed by the service.
There is a `Bootstrap Users` collection that can be run using the `Run collection` tool in Postman that will create 100 users to test the search endpoint with.

The collections will need an environment setup with `scheme`, `port` and `host` variables setup with values of `http`, `8080` and `localhost` respectively.

### Run tests

The service contains a collection of unit and integration tests for testing the various layers of the service. Some of the integration tests use docker to spin up dependencies on demand (ie a postgres db) so just be aware that docker is needed to run the tests.

1. From root of the repo
2. Run `go test ./...`
