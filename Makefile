TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/slackboard bin/slackboard-cli bin/slackboard-log

bundle:
	glide install

bin/slackboard: cmd/slackboard/slackboard.go slackboard/*.go
	go build $(GOFLAGS) -o bin/slackboard cmd/slackboard/slackboard.go

bin/slackboard-cli: cmd/slackboard-cli/slackboard-cli.go slackboard/*.go
	go build $(GOFLAGS) -o bin/slackboard-cli cmd/slackboard-cli/slackboard-cli.go

bin/slackboard-log: cmd/slackboard-log/slackboard-log.go slackboard/*.go
	go build $(GOFLAGS) -o bin/slackboard-log cmd/slackboard-log/slackboard-log.go

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

test:
	go test $(TARGETS_NOVENDOR)

clean:
	rm -rf bin/slackboard*
