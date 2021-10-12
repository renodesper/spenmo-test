BINARY_NAME=spenmo-test
VERSION?=0.0.1
SERVICE_PORT?=8000
DOCKER_REGISTRY?= #if set it should finished by /
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true
BUILD_DIR=build

GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet

GOPATH  := $(shell $(GOCMD) env GOPATH)
AIRPATH := $(GOPATH)/bin/air
GREEN   := $(shell tput -Txterm setaf 2)
YELLOW  := $(shell tput -Txterm setaf 3)
WHITE   := $(shell tput -Txterm setaf 7)
RESET   := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

lint: lint-go lint-dockerfile

lint-dockerfile:
# If Dockerfile is present we lint it.
ifeq ($(shell test -e ./Dockerfile && echo -n yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -i $(CONFIG_OPTION) hadolint/hadolint hadolint --ignore DL3007 --ignore DL3018 $(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
endif

lint-go:
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)

clean:
	rm -fr ./bin
	rm -fr ./build

test:
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -race ./... $(OUTPUT_OPTIONS)

test-unit:
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -short -race ./... $(OUTPUT_OPTIONS)

test-integration:
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -race -run ".Integration" ./... $(OUTPUT_OPTIONS)

vendor:
	$(GOCMD) mod vendor
	$(GOCMD) mod tidy

build:
	mkdir -p $(BUILD_DIR)
	rm -rf $(BUILD_DIR)/*
	CGO_ENABLED=0 GO111MODULE=on $(GOCMD) build -ldflags="-s -w" -mod vendor -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

release:
	@echo 'Not implemented yet'

docker-build:
	docker build --rm --tag $(BINARY_NAME) .

docker-release:
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
#	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
#	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

watch:
	test -s ${AIRPATH} || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(GOPATH)/bin
	${AIRPATH}

help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}build             ${RESET} ${GREEN}Build your project and put the output binary in $(BUILD_DIR)/$(BINARY_NAME)${RESET}"
	@echo "  ${YELLOW}clean             ${RESET} ${GREEN}Remove build related file${RESET}"
	@echo "  ${YELLOW}docker-build      ${RESET} ${GREEN}Use the dockerfile to build the container (name: $(BINARY_NAME))${RESET}"
	@echo "  ${YELLOW}docker-release    ${RESET} ${GREEN}Release the container \"$(DOCKER_REGISTRY)$(BINARY_NAME)\" with tag latest and $(VERSION) ${RESET}"
	@echo "  ${YELLOW}help              ${RESET} ${GREEN}Show this help message${RESET}"
	@echo "  ${YELLOW}lint              ${RESET} ${GREEN}Run all available linters${RESET}"
	@echo "  ${YELLOW}lint-dockerfile   ${RESET} ${GREEN}Lint the Dockerfile using 'hadolint/hadolint'${RESET}"
	@echo "  ${YELLOW}lint-go           ${RESET} ${GREEN}Lint all go files using 'golangci/golangci-lint'${RESET}"
	@echo "  ${YELLOW}test              ${RESET} ${GREEN}Run the tests of the project${RESET}"
	@echo "  ${YELLOW}test-unit         ${RESET} ${GREEN}Run the unit tests of the project${RESET}"
	@echo "  ${YELLOW}test-integration  ${RESET} ${GREEN}Run the integration tests of the project${RESET}"
	@echo "  ${YELLOW}vendor            ${RESET} ${GREEN}Copy all packages needed to support builds and tests into the vendor directory${RESET}"
	@echo "  ${YELLOW}watch             ${RESET} ${GREEN}Run the code with 'cosmtrek/air' to have automatic reload on changes${RESET}"
