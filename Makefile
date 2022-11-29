default: test

test:
	cd calculator
	go test -v

build:
	cd cmd/calculator && go text -v ./...
	go build -o calculator && go build -o calculator