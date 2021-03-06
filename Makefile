export GO111MODULE=on

all: bin/slackboard bin/slackboard-cli bin/slackboard-log

bin/slackboard: cmd/slackboard/slackboard.go slackboard/*.go
	go build $(GOFLAGS) -o bin/slackboard cmd/slackboard/slackboard.go

bin/slackboard-cli: cmd/slackboard-cli/slackboard-cli.go slackboard/*.go
	go build $(GOFLAGS) -o bin/slackboard-cli cmd/slackboard-cli/slackboard-cli.go

bin/slackboard-log: cmd/slackboard-log/slackboard-log.go slackboard/*.go
	go build $(GOFLAGS) -o bin/slackboard-log cmd/slackboard-log/slackboard-log.go

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf bin/slackboard*
