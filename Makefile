.PHONY: build run
.SILENT:

build:
	go build -o ./build/app ./cmd/forum/main.go

run: build
	./build/app