all: build

build:
	go build -o app

test:
	go test ./interfaces

clear:
	go clean --testcache