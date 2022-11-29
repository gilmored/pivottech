default: test

test:
	cd calculator
	go test -v

build:
	cd cmd/calculator && go text -v ./...
	go build cmd/calculator && go build -o calculator