all: bin/slackboard bin/slackboard-cli bin/slackboard-log

bin/slackboard: cmd/slackboard/slackboard.go slackboard/*.go
	GO111MODULE=on go build $(GOFLAGS) -o bin/slackboard cmd/slackboard/slackboard.go

bin/slackboard-cli: cmd/slackboard-cli/slackboard-cli.go slackboard/*.go
	GO111MODULE=on go build $(GOFLAGS) -o bin/slackboard-cli cmd/slackboard-cli/slackboard-cli.go

bin/slackboard-log: cmd/slackboard-log/slackboard-log.go slackboard/*.go
	GO111MODULE=on go build $(GOFLAGS) -o bin/slackboard-log cmd/slackboard-log/slackboard-log.go

fmt:
	go fmt ./...

test:
	GO111MODULE=on go test $(TARGETS_NOVENDOR)

clean:
	rm -rf bin/slackboard*
