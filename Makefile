.DEFAULT_GOAL := run

.PHONY = clean fmt vet build run test

clean:
	go clean -i
fmt:clean
	go fmt ./...
vet:fmt
	go vet ./...
build:vet
	go build -o snip.exe
run:build
	snip.exe

test:
	go test ./...
