-include .env

VERSION := $(shell git describe --tags)
HASH := $(shell git rev-parse --short HEAD)
BUILD := $(shell date +%Y%m%d-%H%M)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOBUILDDIR := $(GOBASE)/cmd/nuntius
# GOFILES := $(wildcard *.go)
# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.version=$(VERSION) -X=main.gitHash=$(HASH) -X main.build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
# install: go-install

## compile: Compile the binary.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nStdout:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean

go-compile: go-get go-build

go-build:
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=on go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOBUILDDIR)

# go-generate:
#	@echo "  >  Generating dependency files..."
#	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

# go-get:
#	@echo "  >  Checking if there is any missing dependencies..."
#	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
