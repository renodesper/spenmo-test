FROM golang:1.17-alpine as builder

# I use Makefile and tput, so I need to install them
RUN apk add --no-cache ncurses make

WORKDIR /go/src/gitlab.com/renodesper/spenmo-test
COPY . .

RUN rm -rf vendor .vendor* \
  && make vendor \
  && make build
RUN ls /go/src/gitlab.com/renodesper/spenmo-test

# Copy into the base image
FROM gcr.io/distroless/static:latest

# Copy the bin file
COPY --from=builder /go/src/gitlab.com/renodesper/spenmo-test/build/spenmo-test /spenmo-test
COPY --from=builder /go/src/gitlab.com/renodesper/spenmo-test/config/env/production.toml /production.toml

ENTRYPOINT ["/spenmo-test", "-config", "./production.toml"]
EXPOSE 8000
