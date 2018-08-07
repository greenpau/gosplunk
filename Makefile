.PHONY: test ctest covdir coverage docs linter qtest clean travis
VERSION:=1.0
GITCOMMIT:=$(shell git describe --dirty --always)
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif

all:
	@mkdir -p ./bin/
	@go build -o ./bin/http-event-collector-client examples/http-event-collector-client.go

linter:
	@golint http-event-collector/client/
	@echo "PASS: golint"

test: covdir linter
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out http-event-collector/client/*

ctest: covdir linter
	@richgo version || go get -u github.com/kyoh86/richgo
	@time richgo test $(VERBOSE) "${TEST}" -coverprofile=.coverage/coverage.out http-event-collector/client/go.*

covdir:
	@mkdir -p .coverage

coverage: covdir
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html

docs:
	@rm -rf .doc/
	@mkdir -p .doc/doc
	@godoc -html github.com/greenpau/gosplunk/http-event-collector/client > .doc/doc/index.html
	@echo "Run to serve docs:"
	@echo "    godoc -goroot .doc/ -html -http \":5000\""

clean:
	@rm -rf .doc
	@rm -rf .coverage

qtest:
	@go test -v -run TestNewClient http-event-collector/client/*

travis:
	@cp examples/.splunk.hec.yaml ~/.splunk.hec.yaml
