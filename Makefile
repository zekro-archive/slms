### NAMES AND LOCS ############################
APPNAME  = slms
PACKAGE  = github.com/zekroTJA/$(APPNAME)
LDPAKAGE = static
CONFIG   = $(CURDIR)/config/private.config.yml
BINPATH  = $(CURDIR)/bin
###############################################

### EXECUTABLES ###############################
GO     = go
DEP    = dep
GOLINT = golint
GREP   = grep
###############################################

# ---------------------------------------------

BIN = $(BINPATH)/$(APPNAME)

TAG        = $(shell git describe --tags)
COMMIT     = $(shell git rev-parse HEAD)


ifneq ($(GOOS),)
	BIN := $(BIN)_$(GOOS)
endif

ifneq ($(GOARCH),)
	BIN := $(BIN)_$(GOARCH)
endif

ifneq ($(TAG),)
	BIN := $(BIN)_$(TAG)
endif

ifeq ($(OS),Windows_NT)
	ifeq ($(GOOS),)
		BIN := $(BIN).exe
	endif
endif

ifeq ($(GOOS),windows)
	BIN := $(BIN).exe
endif


PHONY = _make
_make: deps build cleanup

PHONY += build
build: $(BIN) 

PHONY += deps
deps:
	$(DEP) ensure -v

$(BIN):
	$(GO) build  \
		-v -o $@ -ldflags "\
			-X $(PACKAGE)/$(LDPAKAGE).AppVersion=$(TAG) \
			-X $(PACKAGE)/$(LDPAKAGE).AppCommit=$(COMMIT) \
			-X $(PACKAGE)/$(LDPAKAGE).Release=TRUE" \
		$(CURDIR)/cmd/$(APPNAME)

PHONY += test
test:
	$(GO) test -v -cover ./...

PHONY += lint
lint:
	$(GOLINT) ./... | $(GREP) -v vendor || true

PHONY += run
run:
	$(GO) run -v \
		$(CURDIR)/cmd/$(APPNAME) -c $(CONFIG) -l 5

PHONY += cleanup
cleanup:

PHONY += help
help:
	@echo "Available targets:"
	@echo "  #        - creates binary in ./bin"
	@echo "  cleanup  - tidy up temporary stuff created by build or scripts"
	@echo "  deps     - ensure dependencies are installed"
	@echo "  lint     - run linters (golint)"
	@echo "  run      - debug run app (go run) with test config"
	@echo "  test     - run tests (go test)"
	@echo ""
	@echo "Cross Compiling:"
	@echo "  (env GOOS=linux GOARCH=arm make)"
	@echo ""
	@echo "Use different configs for run:"
	@echo "  make CONF=./myCustomConfig.yml run"
	@echo ""


.PHONY: $(PHONY)