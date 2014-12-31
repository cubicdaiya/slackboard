
all: bin/slackboard

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

bin/slackboard: slackboard.go slackboard/*.go
	gom build -ldflags '-s -w' -o bin/slackboard slackboard.go

fmt:
	go fmt ./...

clean:
	rm -rf bin/slackboard
