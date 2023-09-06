# Grading_API



- This is a service that provides the grading API for the grading system.
- The service is written in Go and uses the following libraries:
	- [goose](github.com/pressly/goose/v3) - for database migrations
    - [openapi-generator](github.com/deepmap/oapi-codegen/cmd/oapi-codegen) - for generating the API code
    - [sqlx](github.com/jmoiron/sqlx) - for database access
    - [testify](github.com/stretchr/testify) - for testing
    - [viper](github.com/spf13/viper) - for configuration
    - [gnomock](github.com/orlangure/gnomock) - for integration testing
    - [slog](golang.org/x/exp/slog) - for logging
## endpoints
  check the [grading.openapi3.yaml](api%2Fv1%2Fgrading.openapi3.yaml) file for the endpoints
## Repository Structure
```
├── Makefile
├── Dockerfile
├── Docker_compose.yaml         # Docker compose file to run the service and its dependencies
├── .golangci.yml                  # GolangCI configuration
├── api
│   ├── v1
│   │   ├── grading.openapi.yaml  # OpenAPI 3.0.3 specification
│   │   ├── gen.go                # Script to generate the API code
│   │   ├── service.config.yaml   # Configuration file for the API
│   │   └── api.gen.go            # Generated API code
├── docs
│   ├── architecture              # Architecture documentation
│   │   └── 0000-record-architecture-decisions.md
│   └── diagrams                    # Design documentation
│       └── overview.wsd 		    # Diagrams in PlantUML format
├── service
│   ├── cmd
│   │   └── grading
│   │       └── main.go           # Entrypoint of the service
│   ├── foundation                # Foundation layer that contains the shared libraries
│   │   ├── http
│   │   │   └── http.go           # HTTP server library that wraps the standard library
│   │   ├── db
│   │   │   └── db.go             # Database library that wraps the standard library
|   ├── shared                    # Shared structs and interfaces between the different layers
│   │   └── errors.go             # local error types
│   │   └── domain
│   │       └── domain.go
│   ├── pkg                        # exportable packages that can be used by other services and defines dependencies
│   │   ├── app
│   │   │   └── app.go
│   ├── tests
│   │   ├── integration
│   │   │   ├── integration.go
│   │   │   └── e2e_test.go
│   │   ├── support
│   │   │   └── client
│   │   │       └── client.go
│   │   │   └── storage
│   │   │       └── sqlt        # SQL test helpers
│   │   │           └── sqlt.go
│   ├── internal                   # internal packages that are not exported and defines dependencies
│   │   ├── api                    # API layer
│   │   │   └── http
│   │   │       └── handler.go
│   │   ├── usecase                  # Use case layer that contains use cases and business logic
│   │   │   ├── usecase.go
│   │   │   └── usecase_mock.go
│   │   ├── storage                # Storage layer that contains the database logic
│   │   │   ├── cache
│   │   │   ├── object
│   │   │   └── postgres
│   │   │       ├── stores.go
│   │   │       ├── read.go
│   │   │       └── write.go
│   │   ├── migration              # Database migrations and scripts powered by https://github.com/pressly/goose
│   │   │   └── postgres
│   │   │       ├── migration.go
│   │   │       ├── *.sql          # SQL scripts to be executed by goose
├── go.mod
├── go.sum
└── README.md 

```

## configuration

to configure the service you can either use the environment variables or the config file

### config file
 set the path to the config file using the `CONFIG_FILE` environment variable
 a sample config file would look like this
 ```yaml
 Host: localhost
 Port: 8080
 LogLevel: debug
 ServiceName: test-service
 IsDryRun: true
 DB:
   User: testuser
   Password: testpassword
   Host: localhost
   Port: 5432
   DBName: testdb`
```


### environment variables


| Name        | Description | Default Value |
|-------------|-------------|---------------|
| PORT        | Port to run the service on | 8080 |
| LOGLEVEL    | Log level | info |
| DB.HOST     | Database host | localhost |
| DB.PORT     | Database port | 5432 |
| DB.USER     | Database user | postgres |
| DB.PASSWORD | Database password | postgres |
| DB.DBNAME   | Database name | grading |
| DB_.SSLMODE | Database ssl mode | disable |


## Run the service
- To run the service:
```shell
make run
```

## Run Tests
- To run the tests:
```shell
make test
```

## Run Integration Tests
integration tests require docker to be installed and running
- To run the integration tests:
```shell
make integration-test
```



## Design

![component-0.png](docs%2Fdiagrams%2Fcomponent-0.png)


## TODO
- [ ] add more tests
- [ ] add more documentation
- [ ] add more logging
- [ ] add authentication
- [ ] add authorization
- [ ] add tracing
- [ ] add metrics