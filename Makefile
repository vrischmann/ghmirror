VERSION = $(shell cat VERSION)
COMMIT = $(shell git rev-parse HEAD)

default:
	go build --ldflags="-X main.commit=$(COMMIT) -X main.version=$(VERSION)" -o ghmirror

linux:
	GOOS=linux GOARCH=amd64 go build --ldflags="-X main.commit=$(COMMIT) -X main.version=$(VERSION)" -o ghmirror
