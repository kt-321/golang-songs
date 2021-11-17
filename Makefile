all: build

build:
	go build -o app

test:
	go test ./interfaces
	go test ./queries/userQuery

clear:
	go clean --testcache