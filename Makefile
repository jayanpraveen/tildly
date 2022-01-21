.PHONY: clean test build tidy tildly
BINARY_NAME=tildly

all: test build tildly

tidy: go.mod
	go mod tidy -v

test:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open ./coverage.html 

build: main.go
	go build -o bin/${BINARY_NAME} main.go

tildly:
	./bin/${BINARY_NAME}

darwin:
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin main.go

linux:
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go

windows:
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows main.go

build_all: darwin linux windows

clean:
	go clean
	rm -f ./bin/tildly
	rm -f ./bin/tildly-darwin
	rm -f ./bin/tildly-linux
	rm -f ./bin/tildly-windows
	rmdir bin
	
redis: redis.conf
	redis-server ./redis.conf
