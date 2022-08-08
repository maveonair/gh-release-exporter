.PHONY:  build clean dev test release

VERSION=0.2.1

default: build

build: clean
	CGO_ENABLED=0 go build -o ./dist/gh-release-exporter -a -ldflags '-s' -installsuffix cgo cmd/gh-release-exporter/main.go

clean:
	rm -rf ./dist/*

dev:
	gow run cmd/gh-release-exporter/main.go

test:
	go test -v ./...

release: clean
	goreleaser release
