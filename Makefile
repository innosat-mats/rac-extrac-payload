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
GOLINT := $(GOBIN)/golint
GOTEST := $(GOCMD) test
GOINSTALL := $(GOCMD) install $(GOFLAGS)
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
TOOLS := golang.org/x/lint/golint

all: build lint test install
build:
	$(GOBUILD) ./...
lint: tools
	${GOVET} ./...
	${GOLINT} -set_exit_status ./...
test:
	$(GOTEST) -race ./...
install:
	$(GOINSTALL) -i ./...
clean:
	$(GOCLEAN) -i -cache ./...
tools:
	$(GOGET) $(TOOLS)
