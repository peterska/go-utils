.PHONY: all goutils deps test modtidy gofmt vet modlist

GO111MODULE=on

all: goutils

goutils:
	go build

deps:
	go get -d ./...

test:
	go test 

modtidy:
	go mod tidy

gofmt:
	go fmt

vet:
	go vet

modlist:
	go list -m all

