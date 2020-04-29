# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BIN_DIR=bin


all: build install
build:
		$(GOBUILD) ./...
test:
		$(GOTEST)  ./...
clean:
		$(GOCLEAN) -i ./...
install:
		$(GOINSTALL) -i ./...
