# commands

.PHONY: run build clean run-build

run:
	go run main.go

run-build: clean build
	./dist/chmod-cli

build: clean
	go build -o dist/chmod-cli

install:
	go install

clean:
	rm -rf ./dist