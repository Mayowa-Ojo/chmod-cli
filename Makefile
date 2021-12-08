# commands

.PHONY: run build clean run-build

run:
	go run main.go

run-build: clean build
	./dist/chmod-cli

build:
	go build -o dist/chmod-cli

clean:
	rm -rf ./dist