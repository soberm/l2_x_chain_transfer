PROGRAMS=operator setup
BUILD_DIR=bin

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run

.PHONY: all
all: test build

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: build
build: $(PROGRAMS)

$(PROGRAMS):
	$(GOBUILD) -o $(GOBIN)/$@ -v ./cmd/$@

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(GOBIN)/*

.PHONY: run
run:
	@if [ -z $(program) ]; then \
		echo "Usage: make run program=<program_name>"; \
	else \
		$(GORUN) ./cmd/${program}/*.go; \
	fi
