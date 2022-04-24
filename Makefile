all: build

PHONY: build
build:
	go build -o bin/simplebackup main.go
