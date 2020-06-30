VERSION := "0.1"
PROJECTNAME := "lxc-go-http-api"
BIN_INSTALL_DIR := "/usr/local/bin"
DOCS_INSTALL_DIR := "/usr/share/doc/${PROJECTNAME}/swagger"


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

.PHONY: build
## build: Build the binary
# build: go-clean go-get go-build
build: go-clean go-build


.PHONY: clean
## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean

.PHONY: install
## install: Install binary and docs
install:  Install-bin install-docs

.PHONY: install-bin
## install-bin: Copy binary
install-bin:
	@echo "  >  Copy binary to ${BIN_INSTALL_DIR}/$(PROJECTNAME)"
	@cp $(GOBIN)/$(PROJECTNAME) ${BIN_INSTALL_DIR}
	@chmod 0755 ${BIN_INSTALL_DIR}/$(PROJECTNAME)
	@chown root:root ${BIN_INSTALL_DIR}/$(PROJECTNAME)

.PHONY: install-docs
## install-docs: Copy documentation
install-docs:
	@echo "  >  Copy documentation to ${DOCS_INSTALL_DIR}/swagger.json"
	@test -d ${DOCS_INSTALL_DIR} || mkdir -p ${DOCS_INSTALL_DIR}
	@cp ${SWAGGER_SPEC} ${DOCS_INSTALL_DIR}/swagger.json
	@chmod 0644 ${DOCS_INSTALL_DIR}/swagger.json
	@chown root:root ${DOCS_INSTALL_DIR}/swagger.json

.PHONY: docs
## docs: Generate and validate swagger docs
docs: swagger-generate swagger-validate

.PHONY: swagger-generate
swagger-generate:
	@echo "  >  Building swagger specs"
	@test -d ${DOCS_DIR} || mkdir ${DOCS_DIR}
	@swagger generate spec -m -o './$(SWAGGER_SPEC)'

.PHONY: swagger-validate
swagger-validate:
	@echo "  >  Validate swagger specs"
	@swagger validate './$(SWAGGER_SPEC)'

## all: Build and copy binary
all: build install

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
