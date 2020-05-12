# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install


all: build test install
build:
		$(GOBUILD) ./...
test:
		$(GOTEST) ./...
install:
		$(GOINSTALL) -i ./...
clean:
		$(GOCLEAN) -i -cache ./...
