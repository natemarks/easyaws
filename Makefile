.DEFAULT_GOAL := help

# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))
PKG := github.com/natemarks/easyaws
VERSION := 0.0.0
DEFAULT_BRANCH := main
COMMIT := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
INTEGRATION_SCRIPTS := $(shell find ./scripts -type f -name "test_*.sh")


help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

all: run

clean-venv: ## re-create virtual env
	rm -rf .venv
	python3 -m venv .venv
	( \
       source .venv/bin/activate; \
       pip install --upgrade pip setuptools; \
    )

build: ## build the binaries with commit IDs
	mkdir -p build/$(COMMIT)/linux/amd64
	env GOOS=linux GOARCH=amd64 \
	go build  -v -o build/$(COMMIT)/linux/amd64/${OUT} \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${COMMIT}" ${PKG}
	mkdir -p build/$(COMMIT)/darwin/amd64
	env GOOS=darwin GOARCH=amd64 \
	go build  -v -o build/$(COMMIT)/darwin/amd64/${OUT} \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${COMMIT}" ${PKG}

release:  ## Build release versions
	mkdir -p build/$(VERSION)
	env GOOS=linux GOARCH=amd64 \
	go build  -v -o build/$(VERSION)/${OUT}_linux_amd64 \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${VERSION}" ${PKG}
	env GOOS=darwin GOARCH=amd64 \
	go build  -v -o build/$(VERSION)/${OUT}_darwin_amd_64 \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${VERSION}" ${PKG}

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

static: vet lint test

run: server
	./${OUT}

shellcheck: ## Run static code checks
	@echo Run shellcheck against scripts/
	shellcheck scripts/*.sh

clean:
	-@rm ${OUT} ${OUT}-v*

create_test_iam_user:

i_test: create_test_iam_user ## run all of the integration tests
	bash $(INTEGRATION_SCRIPTS)

bump: clean-venv  ## bump version in main branch
ifeq ($(CURRENT_BRANCH), $(MAIN_BRANCH))
	( \
	   source .venv/bin/activate; \
	   pip install bump2version; \
	   bump2version $(part); \
	)
else
	@echo "UNABLE TO BUMP - not on Main branch"
	$(info Current Branch: $(CURRENT_BRANCH), main: $(MAIN_BRANCH))
endif


.PHONY: run build release static upload vet lint shellcheck