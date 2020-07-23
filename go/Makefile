# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
BINDIR=./bin

all: build 

build: deps
	$(GOBUILD) -v -o ${BINDIR}/oslabs github.com/reo-public-tools/openstack-labs/oslabs

clean: cleandeps
	$(GOCLEAN)
	rm -f ${BINDIR}/golangci-lint
	rm -f ${BINDIR}/oslabs

test:
	$(GOTEST) -v ./...


###############
# Dependencies
###############
deps: src/gopkg.in/yaml.v2 src/github.com/gophercloud/gophercloud src/github.com/gophercloud/utils

src/github.com/gophercloud/gophercloud:

	$(GOCMD) get github.com/gophercloud/gophercloud

src/github.com/gophercloud/utils:

	$(GOCMD) get github.com/gophercloud/utils/openstack/clientconfig
	$(GOCMD) get github.com/gophercloud/utils/openstack/compute/v2/flavors
	$(GOCMD) get github.com/gophercloud/utils/openstack/networking/v2/networks

src/gopkg.in/yaml.v2:

	$(GOCMD) get gopkg.in/yaml.v2

cleandeps:
	rm -rf ./src/github.com/gophercloud
	rm -rf ./src/gopkg.in/yaml.v2



###########
# Linting 
###########
./bin/golangci-lint:

	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b ${GOPATH}/bin v1.21.0

lint: ./bin/golangci-lint
      
	./bin/golangci-lint run \
	    -e 'Command.Parse' \
	    src/github.com/reo-public-tools/openstack-labs/*



###########
# Install
###########

/etc/oslabs/oslabs.yml: /etc/oslabs

	cp example_configs/oslabs.yml /etc/oslabs/oslabs.yml

/etc/oslabs:

	mkdir /etc/oslabs
        chown root:root /etc/oslabs
	chmod 744 /etc/oslabs

/usr/bin/oslabs:

	cp ./bin/oslabs /usr/bin/oslabs
	chown root:root /usr/bin/oslabs
	chmod 755 /usr/bin/oslabs

install: build /etc/oslabs/oslabs.yml /usr/bin/oslabs

cleaninstall:
	rm -f /usr/bin/oslabs
	rm -rf /etc/oslabs
