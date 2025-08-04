<div align="center">
 <a href="https://www.speakeasy.com/" target="_blank">
  <img width="1500" height="500" alt="Speakeasy" src="https://github.com/user-attachments/assets/0e56055b-02a3-4476-9130-4be299e5a39c" />
 </a>
 <br />
 <br />
  <div>
   <a href="https://speakeasy.com/docs/create-client-sdks/" target="_blank"><b>Docs Quickstart</b></a>&nbsp;&nbsp;//&nbsp;&nbsp;<a href="https://go.speakeasy.com/slack" target="_blank"><b>Join us on Slack</b></a>
  </div>
 <br />

</div>

# The RESTful API Template Project [Golang]

## How To Use This Repo

This repo is intended to be used by Golang developers seeking to understand the building blocks of a simple and well-constructed REST API service. We have built a simple CRUD API which exhibits the characteristics we expect our own developers to apply to the APIs we build at Speakeasy:

- **Entity-based**: The resources available should represent the domain model. Each resource should have the CRUD methods implemented (even if not all available to API consumers). In our template, we have a single resource defined (users.go). However other resources could be easily added by copying the template and changing the logic of the service layer.
- **Properly Abstracted**: The Transport, service, and data layers are all cleanly abstracted from one another. This makes it easy to make apply updates to the API endpoints
- **Consistent**: It's important that consumers of a service have guaranteed consistency across the entire range of API endpoints and methods. In this service, responses are consistently formatted whether successfully returning a JSON object or responding with an error code. All the service's methods use shared response (http.go) and error (errors.go) handler functions to ensure consistency.
- **Tested**: We believe that a blend of unit and integration testing is important for ensuring that the service maintains its contract with consumers. The service repo therefore contains a collection of unit and integration tests for the various layers of the service.
- **Explorable**: It is important for developers to be able to play with an endpoint in order to understand it. We have provided Postman collections for testing out the REST endpoints exposed by the service. That's why there is a `Bootstrap Users` collection that can be run using the `Run collection` tool in Postman that will create 100 users to test the search endpoint with.

This repo can serve as an educational tool, or be used as a foundation upon which developers can build their own basic API scaffolding to turn API development into a consistent and marignally easier activity.

## Getting Started

### Prerequisites

- Go 1.18 (should still be backwards compatible with earlier versions)

### Running locally

1. From root of the repo
2. Run `docker-compose up` will start the dependencies and server on port 8080

### Running via docker

1. From root of the repo
2. Run `docker-compose up` will start the dependencies and server on port 8080

### Postman

The collections will need an environment setup with `scheme`, `port` and `host` variables setup with values of `http`, `8080` and `localhost` respectively.

### Run tests

Some of the integration tests use docker to spin up dependencies on demand (ie a postgres db) so just be aware that docker is needed to run the tests.

1. From root of the repo
2. Run `go test ./...`
