BINARY_NAME=tildly

all: test build tildly

help :
	@echo "run         : runs the app in port 8080"
	@echo "test        : runs tests & generaters report."
	@echo "build       : builds the binary"
	@echo "etcd        : starts etcd cluster"
	@echo "etcd-list   : displays the available etcd nodes"


tidy: go.mod
	go mod tidy -v

run:
	go run main.go -port=8080

test:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func coverage.out
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

etcd-list:
	etcdctl --write-out=table --endpoints=localhost:2379 member list

etcd: Procfile 
	cd etcd
	echo "etcd: starting local multi-member cluster with 3 nodes"
	goreman -f Procfile start

.PHONY: clean test build tidy tildly etcd help