# Summary

This project was made especially as a requirement for passing the technical test in Spenmo. It is a simple golang application that allows us to create a team/individual wallet where each wallet may has multiple cards.

In the `docs` folder, you can find the schema of the database and the api collection (I use Insomnia) for the project.

## Constraints and Assumptions

- Client makes a read or write request
- Service does processing, stores data, then returns the results
- Service needs to evolve from serving a small amount of users to a lot of users
- Service has high availability

## System Design Plan

The project is written in golang and uses the following technologies:

- Gokit (a toolkit for microservices)
- PostgreSQL (database)

The uses of Gokit was decided since its design is well suited with modern software design like elegant monolith and microservices. Aside from that, it is also a toolkit that I am familiar with. Since Gokit also support the separation of concerns, the project can be easily divided into multiple layers based on the Service Layer pattern. The project is divided into the following layer groups:

- Transport layer, which is responsible for the communication between the application and the rest of the internet (HTTP, Grpc, etc).
- Endpoint layer, which is responsible in implementing the transport into a business logic, with endpoint we can have two different transport to access a specific business logic.
- Service layer, which is responsible for the business logic.

The implementation of all those layers abide to the Clean Code architecture and Hexagonal architecture.

### Starts As Monolith

For the sake of simplicity, the project is assumed to be a monolith. This means that we can run the project as a single service, either as a static binary in a vps or as a container in docker. We can even put the database in a similar box.

When we need to scale it, we can do vertical scalling by adding more power. Since the resources are limited (a single box), we can easily monitor the metrics of the box to determine bottlenecks.

### Starts As Microservice

Once the users grow, we can deploy each part of the project as a microservice. We can start scaling each service horizontally based on its traffic. Deploying the services as containers will make the scaling easier.

We separate each table on the previous database into a separate database. Each service will own its own database and other services cannot directly access the database. Each service can also be splitted into multiple services based on its task, read or write (CQRS pattern). This will allow us to use a different database between the read (replica) and write (master) services.

Since the project has grown into multiple services, we also need to make sure the services can only be accessed from a single gateway and use a private subnet for everything else. We can use an API Gateway for this.

## Observability Plan

If money is not the problem and we don't have enough resources to setup a monitoring system, we can just use a hosted monitoring system like Datadog or New Relic.

We can also use a local monitoring system like Prometheus or InfluxDB. Based on the comparison of the monitoring system on the Prometheus [page](https://prometheus.io/docs/introduction/comparison), I will choose Prometheus instead of InfluxDB. The reason is that I prefer to use another tool for logging like ELK instead of using InfluxDB for both use cases. Quoted from the Prometheus page, Prometheus has "more powerful query language, alerting, and notification functionality". It also has "higher availability and uptime for graphing and alerting".

Based on above use case, I already included a `requestId` metadata on each request that comes in and comes out into/from the service. We can just use it to trace and/or identify the request between all related microservices.

> An example of prometheus implementation can be seen by accessing `"/metrics"` endpoint.

### Relational Tables

![Relational Tables](docs/diagram/spenmo-wallet.png)

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
  build              Build your project and put the output binary in build/spenmo-test
  clean              Remove build related file
  docker-build       Use the dockerfile to build the container (name: spenmo-test)
  docker-release     Release the container "spenmo-test" with tag latest and 0.0.1
  help               Show this help message
  lint               Run all available linters
  lint-dockerfile    Lint the Dockerfile using 'hadolint/hadolint'
  lint-go            Lint all go files using 'golangci/golangci-lint'
  test               Run the tests of the project
  test-unit          Run the unit tests of the project
  test-integration   Run the integration tests of the project
  vendor             Copy all packages needed to support builds and tests into the vendor directory
  watch              Run the code with 'cosmtrek/air' to have automatic reload on changes
```

## Testing

Unfortunately, in this current version, I am only able to finish the integration part of the testing that will test the `service` part of the project. That's why, the only working Make target is `test-integration`. Since the only working test is the integration test, I didn't include the test in the Gitlab CI yaml file (it needs a working PostgreSQL). Below is the command to trigger the test and the current status of the test:

```sh
❯ make test-integration
go test -race -run ".Integration" ./...
?       gitlab.com/renodesper/spenmo-test/cmd                             [no test files]
?       gitlab.com/renodesper/spenmo-test/endpoint                        [no test files]
?       gitlab.com/renodesper/spenmo-test/middleware                      [no test files]
?       gitlab.com/renodesper/spenmo-test/middleware/recover              [no test files]
?       gitlab.com/renodesper/spenmo-test/repository                      [no test files]
?       gitlab.com/renodesper/spenmo-test/repository/postgre              [no test files]
ok      gitlab.com/renodesper/spenmo-test/service               (cached)
ok      gitlab.com/renodesper/spenmo-test/transport/http        (cached)  [no tests to run]
?       gitlab.com/renodesper/spenmo-test/util/card                       [no test files]
?       gitlab.com/renodesper/spenmo-test/util/constant                   [no test files]
?       gitlab.com/renodesper/spenmo-test/util/ctx                        [no test files]
?       gitlab.com/renodesper/spenmo-test/util/error                      [no test files]
?       gitlab.com/renodesper/spenmo-test/util/errors                     [no test files]
?       gitlab.com/renodesper/spenmo-test/util/logger                     [no test files]
?       gitlab.com/renodesper/spenmo-test/util/logger/noop                [no test files]
?       gitlab.com/renodesper/spenmo-test/util/logger/zap                 [no test files]
?       gitlab.com/renodesper/spenmo-test/util/response                   [no test files]
```