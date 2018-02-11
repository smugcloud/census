.PHONY: linux darwin clean all

linux:
	env GOOS=linux GOARCH=amd64 go build -o builds/linux/census

darwin:
	env GOOS=darwin GOARCH=amd64 go build -o builds/darwin/census

all: linux darwin

clean:
	rm -rf builds/