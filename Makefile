install: 
	go install
	export GOROOT=$HOME/go
	export PATH=$PATH:$GOROOT/bin

run: 
	go run .

test:
	go test -v ./src

coverage:
	go test -coverprofile=coverage.out ./src
	go tool cover -html=coverage.out

lint:
	go fmt ./...
	govulncheck ./...

