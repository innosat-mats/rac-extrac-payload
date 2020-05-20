GOPATH ?= $(HOME)/go
GOBIN ?= $(GOPATH)/bin
GOCMD := go
GOBUILD := $(GOCMD) build
GOVET := ${GOCMD} vet
GOLINT := $(GOBIN)/golint
GOTEST := $(GOCMD) test
GOINSTALL := $(GOCMD) install
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get 
TOOLS := golang.org/x/lint/golint

all: build lint test install
build:
	$(GOBUILD) ./...
lint: tools
	${GOVET} ./...
	${GOLINT} ./...
test:
	$(GOTEST) ./...
install:
	$(GOINSTALL) -i ./...
clean:
	$(GOCLEAN) -i -cache ./...
tools:
	$(GOGET) $(TOOLS)
