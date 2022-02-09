# Variables -------------------------------------------------------------------------------------------------------------

APPNAME     = vsummary
REPO        = github.com/gbolo/vsummary
SERVERPKG   = $(REPO)/cmd/vsummary-server
POLLERPKG   = $(REPO)/cmd/vsummary-poller
METAPKG     = $(REPO)/common
INTTESTPKG  = $(REPO)/integrationtest
DATE       ?= $(shell date +%FT%T%z)
VERSION     = 1.2.1-rc1
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
LDFLAGS     = -X $(METAPKG).Version=$(VERSION) -X $(METAPKG).BuildDate=$(DATE) -X $(METAPKG).CommitSHA=$(COMMIT_SHA)
PKGS        = $(shell $(GO) list ./... | grep -v $(INTTESTPKG))
TESTPKGS    = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN         = $(CURDIR)/bin
GO          = go
GOFMT       = gofmt

V ?= 0
Q  = $(if $(filter 1,$V),,@)
M  = $(shell printf "\033[34;1m▶\033[0m")

# Targets to build Go tools --------------------------------------------------------------------------------------------

$(BIN):
	@mkdir -p $@

$(BIN)/%: | $(BIN) ; $(info $(M) building $(REPOSITORY)...)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(REPOSITORY) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: REPOSITORY=golang.org/x/tools/cmd/goimports

GOLANGCI_LINT = $(BIN)/golangci-lint
$(BIN)/golangci-lint: REPOSITORY=github.com/golangci/golangci-lint/cmd/golangci-lint

# Targets for our app --------------------------------------------------------------------------------------------------

.PHONY: all
all: $(BIN) server poller;                                        @ ## Build server and poller binaries

.PHONY: server
server: ; $(info $(M) building server executable...)              @ ## Build server binary
	$Q $(GO) build -ldflags '$(LDFLAGS)' -o $(BIN)/$(APPNAME)-server $(SERVERPKG)

.PHONY: poller
poller: ; $(info $(M) building poller executable...)              @ ## Build poller binary
	$Q $(GO) build -ldflags '$(LDFLAGS)' -o $(BIN)/$(APPNAME)-poller $(POLLERPKG)

.PHONY: docker
docker: clean ; $(info $(M) building docker image...)             @ ## Build docker image
	$Q docker build -t gbolo/$(APPNAME):$(VERSION) .

.PHONY: fmt
fmt: ; $(info $(M) running gofmt...)                              @ ## Run gofmt on all source files
	$Q $(GOFMT) -s -l -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: goimports
goimports: | $(GOIMPORTS) ; $(info $(M) running goimports...)     @ ## Run goimports on all source files
	$Q $(GOIMPORTS) -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint
lint: | $(GOLANGCI_LINT) ; $(info $(M) running golangci-lint...)  @ ## Run golangci-lint for code issues
	$Q $(GOLANGCI_LINT) run

.PHONY: unit-test
unit-test: ; $(info $(M) running unit tests...)                   @ ## Run go unit tests
	$Q $(GO) test -v -cover $(TESTPKGS)

.PHONY: integration-test
integration-test: ; $(info $(M) running integration tests...)     @ ## Run integration tests
	$Q $(GO) test -v -cover \
	-coverpkg $(REPO)/db,$(REPO)/poller \
	$(INTTESTPKG)

.PHONY: setup-integration-prereqs
setup-integration-prereqs: ; $(info $(M) setup integration...)    @ ## Setup integration prerequisites
	$Q testdata/scripts/setup-integration-prereqs.sh

.PHONY: down-integration-prereqs
down-integration-prereqs: ; $(info $(M) teardown integration...)  @ ## Shutdown integration prerequisites
	$Q testdata/scripts/setup-integration-prereqs.sh down

.PHONY: vcsim
vcsim: ; $(info $(M) starting vcsim...)                           @ ## Start local vCenter simulator
	testdata/scripts/vcsim.sh

.PHONY: clean
clean: ; $(info $(M) cleaning...)                                 @ ## Cleanup everything
	$Q rm -rvf $(BIN)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: versions
versions:
	@echo "app version $(VERSION) $(COMMIT_SHA)"; $(GO) version
