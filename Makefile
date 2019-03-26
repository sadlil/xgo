# Makefile includes some useful commands to build or format incentives
# More commands could be added

# Variables
PROJECT = xgo
REPO_ROOT = ${GOPATH}/src/github.com/sadlil
ROOT = ${REPO_ROOT}/${PROJECT}

LINTER_PKGS = $(shell glide nv)
LINTER_EXCLUDE = "(^|/)z_.*\.go|(^|/)mocks/|(^|/)mock_.*\.go|(^|/)(_)?resources/|(^|/)vendor/|(^|/)databases/|(^|/)pb/"

fmt:
	@goimports -w pkg
	@gofmt -s -w  pkg

compile: fmt test
	@go install ./pkg/...

# A user can invoke tests in different ways:
#  - make test runs all tests;
#  - make test TEST_TIMEOUT=10 runs all tests with a timeout of 10 seconds;
#  - make test TEST_PKG=./model/... only runs tests for the model package;
#  - make test TEST_ARGS="-v -short" runs tests with the specified arguments;
#  - make test-race runs tests with race detector enabled.
TEST_TIMEOUT = 60
TEST_PKGS ?= ./...
TEST_TARGETS := test-short test-verbose test-race test-cover
.PHONY: $(TEST_TARGETS) test tests
test-short:   TEST_ARGS=-short
test-verbose: TEST_ARGS=-v
test-race:    TEST_ARGS=-race
test-cover:   TEST_ARGS=-cover
$(TEST_TARGETS): test

test:
	@go test -timeout $(TEST_TIMEOUT)s $(TEST_ARGS) $(TEST_PKGS)

clean:
	@rm -rf bin
	@go clean
