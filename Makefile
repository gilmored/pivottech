default: test

test:
	cd calculator
	go test -v

build:
	cd cmd/calculator && go test -v ./...
	go build cmd/calculator && go build -o calculator