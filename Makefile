# Go parameters
GOPATH ?= $(HOME)/go
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOINSTALL := $(GOCMD) install

all: build test install
build:
	$(GOBUILD) -ldflags='-extldflags=-static' ./...
test:
	$(GOTEST) -ldflags='-extldflags=-static' ./...
install:
	$(GOINSTALL) -i -ldflags='-extldflags=-static' ./...
clean:
	$(GOCLEAN) -i -cache ./...
dist: $(distfiles)

