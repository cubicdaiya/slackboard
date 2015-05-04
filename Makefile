
all: bin/slackboard bin/slackboard-cli bin/slackboard-log

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

bin/slackboard: slackboard.go slackboard/*.go
	gom build $(GOFLAGS) -o bin/slackboard slackboard.go

bin/slackboard-cli: slackboard-cli.go slackboard/*.go
	gom build $(GOFLAGS) -o bin/slackboard-cli slackboard-cli.go

bin/slackboard-log: slackboard-log.go slackboard/*.go
	gom build $(GOFLAGS) -o bin/slackboard-log slackboard-log.go

fmt:
	go fmt ./...

clean:
	rm -rf bin/slackboard*
