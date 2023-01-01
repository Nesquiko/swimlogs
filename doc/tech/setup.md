# Project setup

## Used tools

### OpenApi

OpenApi is a standard format for defining structure of an API. For more information
see [OpenApi Specification](https://swagger.io/specification/).

OAS (OpenApi Specification) for this project are located in _./doc/openapi_.
The names of schemas start with semantic version (e.g. `v1.0.0_swimlogsApi.yaml`),
which is equal to the version attribute in the schema.

### OpenApi Generator

[Deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen) is used for
generating boilerplate code from a OAS.

Install it with command:
`go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest`

Configuration files for generating server code is in _./generator/config_ and
generated code will be placed in _./generator/generated_.

To generate code use defined **openapi-gen** task with [Task tool](#task).

### Task

[Task](https://taskfile.dev/usage/) is a task runner / build tool.

Install it with command:
`go install github.com/go-task/task/v3/cmd/task@latest`

Tasks for this project are defined in **Taskfile.yaml** file.

### Migrate

SQL migrations are done with [Migrate](https://github.com/golang-migrate/migrate).

For installation refer to its github page.

SQL schemas are in _./migrations_ directory.
