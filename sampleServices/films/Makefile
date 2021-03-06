#-include .env

VERSION := 1.0.0
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
PROJDIR := $(shell pwd)
TARGET := $(PROJECTNAME)

K8SRUNNING := $(shell dev/microk8sStatus.sh)
export DOCKER_IP = `(dev/getdockerip.sh)`

ASSETS := $(PROJDIR)/assets/images
ARTIFACTS := $(PROJDIR)/artifacts
BUILDS := $(PROJDIR)/builds
DOCS := $(PROJDIR)/docs
LOGS := $(PROJDIR)/logs
GITTEST := .git/HEAD
FOSSATEST := .fossa.yml
PREFLIGHT := .pr_preflight_check

# Go related variables.
GOBASE := $(shell cd ../../;pwd)
GOPATH := $(GOBASE)
export GOPATH = $(GOBASE)

#unexport Go modules needed if it is set true globally
unexport GO111MODULE

GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOLINT := $(shell which golint)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
GOCOVERAGE := $(ARTIFACTS)/coverage.out
GOLINTREPORT := $(ARTIFACTS)/lint.out
GOSECREPORT := $(ARTIFACTS)/gosec.out
GOVETREPORT := $(ARTIFACTS)/govet.out
GOTESTREPORT := https://sonarcloud.io/dashboard?id=PavedRoad_films

GIT_TAG := $(shell git describe 2> /dev/null)

SHELL := /bin/bash

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.GitTag=$(GIT_TAG)"

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

.PHONY: check build compile sonar-scanner

all: $(PREFLIGHT) $(FOSSATEST) compile check

## compile: Compile the binary.
compile: $(LOGS) $(ARTIFACTS) $(ASSETS) $(DOCS) $(BUILDS) api-doc
	@echo "  Compiling"
	@-$(MAKE) -s build

## clean: Remove dep, vendor, binary(s), and executes go clean
clean:
	@echo "  execute go-clean"
	@-rm $(GOBIN)/$(PROJECTNAME)* 2> /dev/null || true
	@-rm -R vendor Gopkt.* 2> /dev/null || true
	@-$(MAKE) go-clean

## build: Build the binary for linux / mac x86 and amd
build: go-get
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(GOFILES)
# make this conditional on build GOARCH
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS="darwin" GOARCH="amd64" go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(GOFILES)
	cp $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH)
	cp $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(PROJECTNAME)
	cp $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(BUILDS)/$(PROJECTNAME)-"darwin"-"amd64"


## deploy: Deploy image to repository and k8s cluster
deploy:
ifeq "$(K8SRUNNING)" "down"
	@echo "  >  Starting k8s for deployment..."
	dev/microk8sStart.sh
endif
	@echo "  >  Starting k8s is up..."
	@dev/kube-config.sh
#       wait for the registry service to be ready	
	@echo "  >  Wait for registry to come up..."
	@sleep 20
	@echo "  >  Build image and deploy..."
	@skaffold run -f manifests/skaffold.yaml

Gopkg.toml:
	@echo "  >  initialize dep support..."
	$(shell (export GOPATH=$(GOPATH);dep init))

go-get: Gopkg.toml get-deps $(ASSETS)
	@echo "  >  Creating dependencies graph png..."
	$(shell (export GOPATH=$(GOPATH);dep status -dot | dot -T png -o $(ASSETS)/$(PROJECTNAME).png))

api-doc: $(DOCS)
	@echo "  >  Generate swagger specification..."
	$(shell (export GOPATH=$(GOPATH);swagger generate spec -m -t $(PROJECTNAME) -o docs/api.json))
	$(shell (grep -q '"definitions": {' docs/api.json || sed -i -e 's/"swagger": "2.0",/"swagger": "2.0",\n"definitions": { "Streams": { "type": "object", "allOf": [ { "properties": {} } ] } },/' docs/api.json))
	@echo "  >  Generate HTML..."
	pretty-swag -i docs/api.json -o docs/api.html
	@echo "  >  Done"

get-deps:
	@echo "  >  dep ensure..."
	$(shell (GOPATH=$(GOPATH);dep ensure $?))

## install: Install packages or main
install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

## check: Start services and execute static code analysis and tests
check: lint sonar-scanner $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)
	@echo "  >  Starting cockroach DB..."
	docker-compose -f manifests/docker-db-only.yaml up -d
	@echo "  >  Waiting cockroach DB is ready..."
	@sleep 5
	@echo "  >  running to tests..."
	go test -coverprofile=$(GOCOVERAGE) -v ./...
	@echo "  >  Stopping cockroach DB..."
	docker-compose -f manifests/docker-db-only.yaml down

sonar-scanner: $(ARTIFACTS)
	./sonarcloud.sh

## show-coverage: Show go code coverage in browser
show-coverage:
	go tool cover -html=$(GOCOVERAGE)

## show-test: Show sonarcloud test report
show-test:
	xdg-open $(GOTESTREPORT)

## show-devkit: Show documentation for Devkit
show-devkit:
	xdg-open http://localhost:5000/microk8sDevKit.html


lint: $(GOFILES)
	@echo -n "  >  running lint..."
	@echo $?
	$(GOLINT) $? > $(GOLINTREPORT)
	@echo "  >  running gosec... > $(GOSECREPORT)"
	$(shell (export GOPATH=$(GOPATH);gosec -fmt=sonarqube -tests -out $(GOSECREPORT) -exclude-dir=.templates ./...))
	@echo "  >  running go vet... > $(GOVETREPORT)"
	$(shell (export GOPATH=$(GOPATH);go vet ./... 2> $(GOVETREPORT)))

	@echo "  >  running FOSSA license scan."
	@FOSSA_API_KEY=$(FOSSA_API_KEY) fossa analyze


## fmt: Run gofmt on all code
fmt: $(GOFILES)
	@gofmt -l -w $?

## simplify: Run gofmt with simplify option
simplify: $(GOFILES)
	@gofmt -s -l -w $?

## k8s-start: Start local microk8s server and update configurations
k8s-start:
	@echo "  > dev/microk8sStart.sh"
	dev/microk8sStart.sh

## k8s-stop: Stop local k8s cluster and delete skaffold deployments
k8s-stop:
	skaffold delete -f manifests/skaffold.yaml
	dev/microk8sStop.sh

## k8s-status: Print the status of the local cluster up or down
k8s-status:
	@echo -n "  >  microk8s is "
	@echo $(K8SRUNNING)

## help: Print possible commands
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## docker-build: Build docker images for use with docker-compose
docker-build:
	docker build -f manifests/Dockerfile -t acme-demo/users .
	docker build -f manifests/Dockerfile -t acme-demo/usersinitdb .

## up: Start service using docker-compose
up:
	docker-compose -f manifests/docker-compose.yaml up -d


## down: Stop service using docker-compose
down:
	docker-compose -f manifests/docker-compose.yaml down

$(ASSETS):
	$(shell mkdir -p $(ASSETS))

$(ARTIFACTS):
	$(shell mkdir -p $(ARTIFACTS))

$(BUILDS):
	$(shell mkdir -p $(BUILDS))

$(DOCS):
	$(shell mkdir -p $(DOCS))

$(LOGS):
	$(shell mkdir -p $(LOGS))

# Preflight: ensure all requirements are met
# Once they are the $(PREFLIGH) file will be created
$(PREFLIGHT):
	dev/preflight.sh


$(FOSSATEST):
	fossa init

