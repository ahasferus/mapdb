PROJPATH=$(shell pwd)
PROJNAME=$(shell basename $(PWD))

GOBIN=$(PROJPATH)/bin
GOSOURCE=$(PROJPATH)/cmd

SERVER_NAME=server
CLIENT_NAME=client

.DEFAULT: build

.PHONY: build
build: clean config build-server build-client
	echo "Finished build of $(PROJNAME)"

.PHONY: build-server
build-server:
	cd $(GOSOURCE)/$(SERVER_NAME); \
	go build -o $(GOBIN)/$(SERVER_NAME) -v;

.PHONY: build-client
build-client:
	cd $(GOSOURCE)/$(CLIENT_NAME); \
	go build -o $(GOBIN)/$(CLIENT_NAME) -v;

.PHONY: config
config:
	mkdir $(GOBIN)

.PHONY: clean
clean:
	rm -rf $(GOBIN)

