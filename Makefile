PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/build/bin

.PHONY: build
build:
	@go build -mod=mod -o $(GOBIN)/itvbackend $(GOBASE)/cmd/itvbackend/*.go

.PHONY: migrate-db
migrate-db:
	@migrate -path ./migrations -database postgres://postgres:postgres@localhost:5432/itvbackend?sslmode=disable -verbose up

.PHONY: test
test:
	@go test -mod=mod -race -count 100 $(GOBASE)/internal/app/...

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: clean
clean:
	@rm -fR $(GOBIN)
    @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
