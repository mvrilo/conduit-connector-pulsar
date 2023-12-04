.PHONY: build test lint generate test-pulsar

VERSION=$(shell git describe --tags --dirty --always)

build: generate lint
	go build \
		-ldflags "-X 'github.com/conduitio/conduit-connector-pulsar.version=${VERSION}'" \
		-o conduit-connector-pulsar \
		cmd/connector/main.go

test-pulsar:
	# run required docker containers, execute integration tests, stop containers after tests
	# docker compose -f test/docker-compose-pulsar.yml up --quiet-pull -d --wait
	go test $(GOTEST_FLAGS) -race ./...

test: test-pulsar

generate:
	go generate ./...

lint:
	go vet ./...
	golangci-lint run
