# Notes

Notes is an example application, implementing a note taking service, for the purpose to demonstrating structuring a Go application.

This is not intended to be taken as gospel, but as a collection of patterns I can used in developing Go applications.

## Requirements

  - Go 1.9+
  - [Go Dep](https://github.com/golang/dep)
  - [Migrate](https://github.com/golang-migrate/migrate)


## Getting Started

Setup the database with required tables:

```sh
migrate -database postgres://username:password@host/dbname?sslmode=disable -path migrations up
```

Install Go dependencies:

```sh
dep ensure -v
```

## Running

```sh
go run cmd/main.go
```
