run: 
	go run main.go

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	go fmt ./...
	govulncheck ./...

install: 
	go install
	export GOROOT=$HOME/go
	export PATH=$PATH:$GOROOT/bin
