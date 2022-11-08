GOPATH ?= $(HOME)/go
GOBIN ?= $(GOPATH)/bin
GOCMD := go
ifeq ($(OS),Windows_NT)
SHELL := git-bash.exe
endif
# Fetch build info
HEAD := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
VERSION := $(shell git describe --tags)
# Prepare flags
GOLDFLAGS += -X main.Head=$(HEAD)
GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

GOBUILD := $(GOCMD) build $(GOFLAGS)
GOVET := ${GOCMD} vet
GOTEST := $(GOCMD) test
GOINSTALL := $(GOCMD) install $(GOFLAGS)
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
GOSTATIC := staticcheck
TOOLS := honnef.co/go/tools/cmd/staticcheck@latest

all: build lint test install
build:
	$(GOBUILD) ./...
lint: tools
	${GOVET} ./...
	${GOSTATIC} ./...
test:
	$(GOTEST) -race ./...
install:
	$(GOINSTALL) ./...
clean:
	$(GOCLEAN) -cache ./...
tools:
	$(GOCMD) install $(TOOLS)
