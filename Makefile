TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/slackboard bin/slackboard-cli bin/slackboard-log

bundle:
	glide install

bindata:
	go-bindata -pkg slackboard ./ui/* && mv bindata.go slackboard/

bin/slackboard: cmd/slackboard/slackboard.go slackboard/*.go bindata
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/slackboard cmd/slackboard/slackboard.go

bin/slackboard-cli: cmd/slackboard-cli/slackboard-cli.go slackboard/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/slackboard-cli cmd/slackboard-cli/slackboard-cli.go

bin/slackboard-log: cmd/slackboard-log/slackboard-log.go slackboard/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/slackboard-log cmd/slackboard-log/slackboard-log.go

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf bin/slackboard*
