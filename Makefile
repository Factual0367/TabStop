BINARY_NAME=TabStop

build: linux windows

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe

clean:
	rm -rf bin/$(BINARY_NAME)-*

all: clean build