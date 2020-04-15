# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BIN_DIR=bin


all: test build
build: 
		$(GOBUILD) -o $(BIN_DIR) -v ./...
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
