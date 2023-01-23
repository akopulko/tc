BINARY_NAME=tc

build:
	GOARCH=arm64 GOOS=darwin go build -mod vendor -o bin/${BINARY_NAME} app/*.go

test:
	go test -v ./...

clean:
	go clean
	rm bin/${BINARY_NAME}