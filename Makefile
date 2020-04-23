# Command variables
GO_BUILD 	= go build
GO_PLUGIN 	= go build -buildmode=plugin
INSTALL 	= /usr/bin/install
MKDIR 		= mkdir -p
RM 		= rm
CP 		= cp

# Optimization build processes
#CPUS ?= $(shell nproc)
#MAKEFLAGS += --jobs=$(CPUS)

OS = $(shell uname -s | tr 'A-Z' 'a-z')
ARCH = amd64

ifeq ($(shell uname -m), x86_64)
		ARCH = amd64
endif

# Project variables
PROJECT_PKG ?= github.com/bazarms/jobflow
PROJECT_PATH ?= $(GOPATH)/src/go/$(PROJECT_PKG)
PROJECT_BIN_DIR ?= bin
PROJECT_PLUGIN_DIR ?= plugins

# Compilation variables
PROJECT_BUILD_SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
PROJECT_BUILD_OSARCH = darwin/amd64 linux/amd64
PROJECT_BUILD_PLUGINS = shell github gox
PROJECT_BUILD_TARGET = jobflow

all: clean build plugins

# Build targets multiple platforms
build: clean
	for osarch in $(PROJECT_BUILD_OSARCH); do \
		OS=`echo $$osarch | cut -d"/" -f1`; \
		ARCH=`echo $$osarch | cut -d"/" -f2`; \
		echo "Compiling $(PROJECT_BUILD_TARGET) for "$$OS"_"$$ARCH"..." ; \
		GOOS=$$OS GOARCH=$$ARCH go build -ldflags="-s -w" -o $(PROJECT_BIN_DIR)"/"$$OS"_"$$ARCH"/"$(PROJECT_BUILD_TARGET); \
	done

optimize:
	for osarch in $(PROJECT_BUILD_OSARCH); do \
		OS=`echo $$osarch | cut -d"/" -f1`; \
		ARCH=`echo $$osarch | cut -d"/" -f2`; \
		echo "Optimizing $(PROJECT_BUILD_TARGET) for "$$OS"_"$$ARCH"..." ; \
		upx --brute $(PROJECT_BIN_DIR)"/"$$OS"_"$$ARCH"/"$(PROJECT_BUILD_TARGET) ; \
		for plugin in $(PROJECT_BUILD_PLUGINS); do \
		echo "Optimizing "$$plugin".so for "$$OS"_"$$ARCH"..." ; \
			upx --brute $(PROJECT_BIN_DIR)"/"$$OS"_"$$ARCH"/"$(PROJECT_PLUGIN_DIR)"/"$$plugin".so" ; \
		done; \
	done

test-unit:
	go test -count 1 -v -tags=unit ./...

test-integration:
	go test -count 1 -v -tags=integration ./test/integration/jobflow -inventory data/inventory.yml -plugin-dir ../../../bin/$(OS)_$(ARCH)/plugins -verbosity 0 -args exec data/flow.yml

fmt:
	gofmt -s -l -w $(PROJECT_BUILD_SRCS)

deps:
	go get -u github.com/mitchellh/gox

clean:
	-$(RM) -rf bin

plugins:
	for osarch in $(PROJECT_BUILD_OSARCH); do \
		OS=`echo $$osarch | cut -d"/" -f1`; \
		ARCH=`echo $$osarch | cut -d"/" -f2`; \
		for plugin in $(PROJECT_BUILD_PLUGINS); do \
			echo "Compiling "$$plugin".so for "$$OS"_"$$ARCH"..." ; \
			$(GO_PLUGIN) -o $(PROJECT_BIN_DIR)"/"$$OS"_"$$ARCH"/"$(PROJECT_PLUGIN_DIR)"/"$$plugin".so" $(PROJECT_PKG)/$(PROJECT_PLUGIN_DIR)/$$plugin; \
		done; \
	done

distclean:

install:

.PHONY: test-unit test-integration clean distclean install plugins
