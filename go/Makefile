# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
BINDIR=./bin

all: build 

build:
	$(GOBUILD) -v -o ${BINDIR}/oslabs github.com/reo-public-tools/openstack-labs/oslabs

clean:
	$(GOCLEAN)
	rm -f ${BINDIR}/golangci-lint
	rm -f ${BINDIR}/oslabs

test:
	$(GOTEST) -v ./...


lint-all: install-golangci-lint lint


install-golangci-lint: 

	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b ${GOPATH}/bin v1.21.0


lint:

	./bin/golangci-lint run 
