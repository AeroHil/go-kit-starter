# go-kit-starter
A [go-kit](https://gokit.io/) boiler plate based on some best practices and my opinions

## Introduction
[Go] is a great general-purpose language, but microservices require a certain amount of 
specialized support. RPC safety, system observability, infrastructure integration, 
even program design. Go kit fills in the gaps left by the standard library, 
and makes Go a first-class language for writing microservices in any organization.

[Go kit] is a toolkit for microservices. It provides guidance and solutions for
most of the common operational and infrastructural concerns. Allowing you to
focus your mental energy on your business logic. It provides the building blocks
for separating transports from business domains; making it easy to switch one
transport for the other. Or even service multiple transports for a service at
once.

## Description
I've spent a lot of time looking through GitHub trying to find the best layout for 
microservices pattern made with Go-Kit. I needed a simple and yet somewhat complete 
template to start pushing builds out at the production level. So naturally, I spend 
a lot of time combining what I like the most from each of these fantastic repositories.

- https://github.com/SeamPay/gokita
- https://github.com/antklim/go-microservices
- https://github.com/kujtimiihoxha/kit

### Package - main go-kit components

- Transport layer (grpc/http)
- Endpoint layer
- Service layer

### Middleware

- Endpoint Logging Middleware
- Endpoint Instrumenting Middleware
- Service Logging Middleware
- Service Instrumenting Middleware

### Observability & Tracing

- Prometheus
- Zipkin

### Misc

- Docker Integration (Dockerfile + docker-compose.yml)
- [Testify](https://github.com/stretchr/testify) (A testing toolkit, but easily removable)

## Project Structure

```
.
├── .gitignore
├── README.md
├── build
│   └── Dockerfile
├── client
│   └── greeter_client.go
├── cmd
│   ├── main.go
│   └── service
│       ├── command.go
│       ├── serve.go
│       ├── service.go
│       └── service_util.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── pb
│   ├── compile.sh
│   ├── greeter.pb.go
│   └── greeter.proto
└── pkg
    ├── common
    │   └── constants.go
    ├── endpoint
    │   ├── endpoint.go
    │   └── middleware.go
    ├── service
    │   ├── middleware.go
    │   └── service.go
    ├── test
    │   ├── service_test.go
    │   └── wiring_test.go
    └── transport
        ├── grpc
        │   ├── handler.go
        │   └── transport.go
        └── http
            ├── handler.go
            └── transport.go
```
### Directories/Files
I won't be going over the obvious files, but I do want to address most files for beginners

`build/Dockerfile` - docker file of how this microservice is containerize

`client/greeter_client.go` - gRPC client code to connect and communicate with your gRPC server and service

`cmd/*` - point of entry code [main.go --> command.go --> serve.go --> service.go]

`docker-compose.yml` - if you have multiple services, it would be easy to manage them all here

`pb/*` - protobuf and generated code for gRPC

`pkg/common/*` - constants and other utils methods used for the project

`pkg/endpoint/endpoing.go` - endpoint code mapped to the service

`pkg/endpoint/middleware.go` - endpoint middleware for logging and instrumentation for endpoints

`pkg/service/service.go` - service code for business logic

`pkg/service/middleware.go` - service middleware for logging and instrumentation mapped to the service

`pkg/test/*` - test code that includes service tests and HTTP wiring tests

`pkg/transport/grpc/*` - gRPC handlers and transport code (separated for readability)

`pkg/transport/http/*` - HTTP handlers and transport code (separated for readability) 

As you build upon your microservice, things start to get messy real soon. The project might look spread out 
right now, but it'll be good to grow into. Feel free to move things around as needed.

## Getting Started

### Prerequisite

To get this service up and running you'll need to have a few things installed beforehand:

1. [Go](https://golang.org/doc/install)
2. [Protobuf](https://github.com/google/protobuf)


Install gRPC prerequisite
```
brew install protobuf

and/or

go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Dependencies

The project uses [Go Modules](https://blog.golang.org/using-go-modules), so you will need to use the `go.mod` file to grab all the dependencies like this:
```
go mod vendor
```

If you make any changes during development make sure you use the above command as well as the following to clean up unused dependencies:
```
go mod tidy
```

### Starting Up

#### Docker

Use the following docker-compose to start up the whole service as is

```
docker-compose up --build
```

After you run docker-compose up your services will start up and any change you make to 
your code will automatically rebuild and restart your service.

#### CLI

Use the following to start up the service from the command line

To make a build
```
go build cmd/main.go
```
To run the server after the `main` executable is built
```
./main server -environment=development -debug-addr=:8080 -grpc-addr=:8082 -http-addr=:8081 -zipkin-addr=http://localhost:9411/api/v2/spans
```

