-include .env

VERSION := $(shell git describe --tags)
HASH := $(shell git rev-parse --short HEAD)
BUILD := $(shell date +%Y%m%d-%H%M)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(firstword $(subst :, ,$(shell go env GOPATH)))
GOBIN := $(GOBASE)/bin
GOBUILDDIR := $(GOBASE)/cmd/nuntius
GO111MODULE := on
# Use linker flags to provide version/build settings
LDFLAGS := -ldflags "-X=main.version=$(VERSION) -X=main.gitHash=$(HASH) -X main.build=$(BUILD)"
GOPACKAGES := ./...

env:
	@echo VERSION $(VERSION)
	@echo HASH $(HASH)
	@echo BUILD $(BUILD)

## test: Run tests.
test: go-test

## fmt: Run go fmt.
fmt: go-fmt

## deps: Get dependencies.
deps: go-deps

## build: Build the binary.
build:
	@-$(MAKE) -s go-deps go-build

## tarball: Make a tarball.
tarball:
	@tar -C $(GOBIN) -zcvf bin/nuntius-$(VERSION).linux-amd64.tar.gz nuntius

go-deps:
	@echo ">> Getting dependencies"
ifdef GO111MODULE
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) go mod download
else
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(GOPACKAGES)
endif

go-build:
	@echo ">> Building binary"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOBUILDDIR)

go-fmt:
	@echo ">> Checking code style"
	@fmtRes=$$(gofmt -d $$(find . -path ./vendor -prune -o -name '*.go' -print)); \
	if [ -n "$${fmtRes}" ]; then \
		echo "gofmt checking failed!"; echo "$${fmtRes}"; echo; \
		echo "Please ensure you are using $$(go version) for formatting code."; \
		exit 1; \
	fi

.PHONY: go-test
go-test:
	@echo ">> Running all tests"
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) go test -v -cover $(GOOPTS) $(GOPACKAGES)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
