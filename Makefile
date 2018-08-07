.PHONY: test ctest covdir coverage docs linter qtest
VERSION:=1.0
GITCOMMIT:=$(shell git describe --dirty --always)
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif

all:
	@echo "WARN: please run 'make test'"

linter:
	@golint
	@echo "PASS: golint"

test: covdir linter
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out

ctest: covdir linter
	@richgo version || go get -u github.com/kyoh86/richgo
	@time richgo test $(VERBOSE) "${TEST}" -coverprofile=.coverage/coverage.out

covdir:
	@mkdir -p .coverage

coverage:
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html

docs:
	@mkdir -p .doc
	@godoc -html github.com/greenpau/go-splunk-ec-client > .doc/index.html
	@echo "Run to serve docs:"
	@echo "    godoc -goroot .doc/ -html -http \":5000\""

clean:
	@rm -rf .doc
	@rm -rf .coverage

qtest:
	@go test -v -run TestNewClient
