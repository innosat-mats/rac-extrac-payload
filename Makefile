# Go parameters
GOCMD=go
GOPATH ?= ~/go/bin
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=${GOCMD} vet
GOLINT=$(GOPATH)/golint
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install


all: build lint test install
build:
		$(GOBUILD) ./...
lint:
		${GOVET} ./...
		${GOLINT} ./...
test:
		$(GOTEST) ./...
install:
		$(GOINSTALL) -i ./...
clean:
		$(GOCLEAN) -i -cache ./...
