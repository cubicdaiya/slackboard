TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/slackboard bin/slackboard-cli bin/slackboard-log

bundle:
	glide install

bin/slackboard: cmd/slackboard/slackboard.go slackboard/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/slackboard github.com/cubicdaiya/slackboard/cmd/slackboard

bin/slackboard-cli: cmd/slackboard-cli/slackboard-cli.go slackboard/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/slackboard-cli github.com/cubicdaiya/slackboard/cmd/slackboard-cli

bin/slackboard-log: cmd/slackboard-log/slackboard-log.go slackboard/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/slackboard-log github.com/cubicdaiya/slackboard/cmd/slackboard-log

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf bin/slackboard*
