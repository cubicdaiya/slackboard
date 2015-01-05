
all: bin/slackboard bin/slackboard-cli

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

bin/slackboard: slackboard.go slackboard/*.go
	gom build -ldflags '-s -w' -o bin/slackboard slackboard.go

bin/slackboard-cli: slackboard-cli.go slackboard/*.go
	gom build -ldflags '-s -w' -o bin/slackboard-cli slackboard-cli.go

fmt:
	go fmt ./...

clean:
	rm -rf bin/slackboard*
