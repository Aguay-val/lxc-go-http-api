VERSION := "0.1"
PROJECTNAME := "lxc-go-http-api"
INSTALL_DIR := "/usr/local/bin"

# Go related variables
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Go Swagger related variables
DOCS_DIR := docs
SWAGGER_SPEC := ${DOCS_DIR}/swagger.json

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION)"

.PHONY: compile
## compile: Compile the binary.
# compile: go-clean go-get go-build
compile: go-clean go-build


.PHONY: clean
## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean

.PHONY: install
## install: Copy binary
install:
	@echo "  >  Copy binary to ${INSTALL_DIR}/$(PROJECTNAME)"
	@cp $(GOBIN)/$(PROJECTNAME) ${INSTALL_DIR}
	@chmod 0755 ${INSTALL_DIR}/$(PROJECTNAME)
	@chown root:root ${INSTALL_DIR}/$(PROJECTNAME)

.PHONY: swagger-generate
swagger-generate:
	@test -d ${DOCS_DIR} || mkdir ${DOCS_DIR}
	swagger generate spec -o './$(SWAGGER_SPEC)'

.PHONY: swagger-validate
swagger-validate:
	swagger validate './$(SWAGGER_SPEC)'

## all: Compile and copy binary
all: compile install

go-build:
	@echo "  >  Building binary"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

# go-get:
# 	@echo "  >  Checking if there is any missing dependencies"
# 	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sort |  sed -e 's/^/ /'
