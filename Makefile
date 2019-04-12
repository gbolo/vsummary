# Variables -------------------------------------------------------------------------------------------------------------

APPNAME     = vsummary
REPO        = github.com/gbolo/vsummary
SERVERPKG   = $(REPO)/cmd/vsummary-server
POLLERPKG   = $(REPO)/cmd/vsummary-poller
METAPKG     = $(REPO)/common
DATE       ?= $(shell date +%FT%T%z)
VERSION     = 1.0
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
LDFLAGS     = -X $(METAPKG).Version=$(VERSION) -X $(METAPKG).BuildDate=$(DATE) -X $(METAPKG).CommitSHA=$(COMMIT_SHA)
PKGS        = $(or $(PKG),$(shell $(GO) list ./...))
TESTPKGS    = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN         = $(CURDIR)/bin
GO          = go

V ?= 0
Q  = $(if $(filter 1,$V),,@)
M  = $(shell printf "\033[34;1m▶\033[0m")

# Targets to build Go tools --------------------------------------------------------------------------------------------

$(BIN):
	@mkdir -p $@

$(BIN)/%: | $(BIN) ; $(info $(M) building $(REPOSITORY)...)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOCACHE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(REPOSITORY) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: REPOSITORY=golang.org/x/tools/cmd/goimports

GOLINT = $(BIN)/golint
$(BIN)/golint: REPOSITORY=golang.org/x/lint/golint

# Targets for our app --------------------------------------------------------------------------------------------------

.PHONY: all
all: $(BIN) server poller;                                    @ ## Build server and poller

.PHONY: server
server: ; $(info $(M) building server executable...)          @ ## Build server binary
	$Q $(GO) build -ldflags '$(LDFLAGS)' -o $(BIN)/$(APPNAME)-server $(SERVERPKG)

.PHONY: poller
poller: ; $(info $(M) building poller executable...)          @ ## Build poller binary
	$Q $(GO) build -ldflags '$(LDFLAGS)' -o $(BIN)/$(APPNAME)-poller $(POLLERPKG)

.PHONY: docker
docker: clean ; $(info $(M) building docker image...)         @ ## Build docker image
	$Q docker build -t gbolo/$(APPNAME):$(VERSION) .

.PHONY: fmt
fmt: ; $(info $(M) running gofmt...)                          @ ## Run gofmt on all source files
	$Q $(GO) fmt ./...

.PHONY: goimports
goimports: | $(GOIMPORTS) ; $(info $(M) running goimports...) @ ## Run goimports on all source files
	$Q $(GOIMPORTS) -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint...)            @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: test
test: ; $(info $(M) running go test...)                       @ ## Run go unit tests
	$Q $(GO) test -v -cover $(TESTPKGS)

.PHONY: clean
clean: ; $(info $(M) cleaning...)                             @ ## Cleanup everything
	$Q rm -rvf $(BIN)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
