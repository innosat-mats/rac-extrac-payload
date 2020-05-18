# Go parameters
GOPATH ?= $(HOME)/go
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOVET := ${GOCMD} vet
GOLINT := golint
GOINSTALL := $(GOCMD) install

all: build lint test install
build:
	$(GOBUILD) -ldflags='-extldflags=-static' ./...
lint:
	${GOVET} ./...
	${GOLINT} ./...
test:
	$(GOTEST) -ldflags='-extldflags=-static' ./...
install:
	$(GOINSTALL) -i -ldflags='-extldflags=-static' ./...
clean:
	$(GOCLEAN) -i -cache ./...

