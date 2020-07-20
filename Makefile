 
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
BINARY_NAME=shorturl
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -race -failfast -count=1 ./...

lint:
	golint -set_exit_status `go list ./...` && go vet ./...

.PHONY: test lint
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GORUN) .