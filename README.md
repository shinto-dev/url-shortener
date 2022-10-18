# url-shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/shinto-dev/url-shortener)](https://goreportcard.com/report/github.com/shinto-dev/url-shortener) 

### Migrations
```shell
make migrateup
```

### Test
To run the unit testcases, use the below command:
```shell
make test
```

### Run Service
To run the http service,
```shell
make run
```

### Adding new migrations
```shell
make addmigration file=<migration name>
```

### Open-API documentation: 
```shell
make generate-doc
open docs/open-api/html/index.html
```

### HTTP Request
The http requests are available at [link](./docs/requests)

### Project Structure
```
|_ cmd 
|_ foundation
|   |_ observation
|   |_ web
|_ internal
|   |_ config
|   |_ core
|       |_ shorturl
        |_ test
|   |_ httpservice
|       |_ appcontext
|       |_ handlers
|_ resources
```
**foundation** This package contains the common packages like observation/logging, web etc which can be extracted as small libraries later.

**config** holds all the configs in the service. The package accepts both yaml or environment variable configs.

**internal/core** contains all the entities in the system. Any business logic related to one entity should be within a single package.

**internal/interactor** [Not present in the current service] This is required when we have apis which requires interaction of multiple core entities. This layer acts as a orchestrator.

**internal/httpservice** is responsible for holding packages which are responsible for exposing the rest APIs.

**cmd** package contains the available commands in the project. This includes starting the server, running migration etc.
