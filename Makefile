
all: bin/slackboard bin/slackboard-cli bin/slackboard-log

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

bin/slackboard: cmd/slackboard/main.go slackboard/*.go
	gom build $(GOFLAGS) -o bin/slackboard github.com/cubicdaiya/slackboard/cmd/slackboard

bin/slackboard-cli: cmd/slackboard-cli/main.go slackboard/*.go
	gom build $(GOFLAGS) -o bin/slackboard-cli github.com/cubicdaiya/slackboard/cmd/slackboard-cli

bin/slackboard-log: cmd/slackboard-log/main.go slackboard/*.go
	gom build $(GOFLAGS) -o bin/slackboard-log github.com/cubicdaiya/slackboard/cmd/slackboard-log

fmt:
	go fmt ./...

clean:
	rm -rf bin/slackboard*
