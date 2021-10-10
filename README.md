# Summary

...

## Setup database

Make sure the postgres has been up and running. If it is not, we can run it in the docker:

```sh
docker run --name postgres-svc-spenmo -p 5432:5432 \
    -e POSTGRES_DB=spenmo \
    -e POSTGRES_USER=user \
    -e POSTGRES_PASSWORD=password \
    -v ~/Tmp/postgres_data:/var/lib/postgresql/data \
    -d postgres:13.2-alpine
```

or

```sh
docker start <container_id>
```

Install `golang-migrate`:

```sh
brew install golang-migrate
```

Migrate the current schema:

```sh
export POSTGRESQL_URL='postgres://user:password@127.0.0.1:5432/spenmo?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path config/db/migrations up
```

Common error:

```sh
error: Dirty database version 2. Fix and force version.
```

Solution:

```sh
migrate -database ${POSTGRESQL_URL} -path config/db/migrations force <version - 1>
```

## Run the app

Use Makefile with hot reload enabled ([air](https://github.com/cosmtrek/air) will automatically installed):

```sh
make watch
```

Use Makefile with hot reload disabled:

```sh
make run
```

Run the main go file:

```sh
go run cmd/main.go
```

For other options, we can use `make help`:

```sh
❯ make help

Usage:
  make <target>

Targets:
  build            Build your project and put the output binary in build/spenmo-test
  clean            Remove build related file
  docker-build     Use the dockerfile to build the container (name: spenmo-test)
  docker-release   Release the container "spenmo-test" with tag latest and 0.0.1
  help             Show this help message
  lint             Run all available linters
  lint-dockerfile  Lint the Dockerfile using 'hadolint/hadolint'
  lint-go          Lint all go files using 'golangci/golangci-lint'
  test             Run the tests of the project
  vendor           Copy all packages needed to support builds and tests into the vendor directory
  watch            Run the code with 'cosmtrek/air' to have automatic reload on changes
```
