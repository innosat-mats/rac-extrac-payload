# Go parameters
GOPATH ?= $(HOME)/go
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOINSTALL := $(GOCMD) install
PWD := $(shell pwd)
ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
LIB_DIR := $(ROOT_DIR)/third-party/linux/lib
distfiles := rac_linux.tar.bz rac_windows.tar.bz rac_darwin.tar.bz

all: build test install
build:
	$(GOBUILD) ./...
test:
	LD_LIBRARY_PATH := $(LIB_DIR) $(GOTEST) ./...
install:
	$(GOINSTALL) -i ./...
	cp -r $(LIB_DIR) $(GOPATH)
clean:
	$(GOCLEAN) -i -cache ./...
	rm -rf $(GOPATH)/lib
	rm -rf $(addprefix $(ROOT_DIR)/, $(distfiles))

dist: $(distfiles)

linux:
	go build -ldflags='-extldflags=-static' -o $(TMPDIR)/bin/rac ./cmd/rac

windows:
	CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags='-extldflags=-static' -o $(TMPDIR)/bin/rac.exe ./cmd/rac

darwin :
	CC=x86_64-apple-darwin19-cc GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build  -o $(TMPDIR)/bin/rac ./cmd/rac

rac_windows.tar.bz: TMPDIR := $(shell mktemp -d)
rac_windows.tar.bz: LIB_DIR := $(ROOT_DIR)/third-party/windows/lib
rac_windows.tar.bz: windows
	$(tarfiles)

rac_darwin.tar.bz: TMPDIR := $(shell mktemp -d)
rac_darwin.tar.bz: LIB_DIR := $(ROOT_DIR)/third-party/darwin/lib
rac_darwin.tar.bz: darwin
	cp -r $(LIB_DIR) $(TMPDIR)
	$(tarfiles)

rac_linux.tar.bz: TMPDIR := $(shell mktemp -d)
rac_linux.tar.bz: LIB_DIR := $(ROOT_DIR)/third-party/linux/lib
rac_linux.tar.bz: linux
	$(tarfiles)

define tarfiles =
tar -cjf $(ROOT_DIR)/$@ -C $(TMPDIR) .
rm -rf $(TMPDIR)
endef
