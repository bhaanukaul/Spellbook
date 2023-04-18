

all: test build

test:
	go test
	rm -rf ./spellbook_test/

# build: build-linux build-darwin build-darwin-arm build-win
build: build-darwin-arm

build-linux:
	GOOS=linux GOARCH=amd64 go build \
		-o bin/spellbook-linux *.go

build-darwin-arm:
	GOOS=darwin GOARCH=arm64 go build \
		-o bin/spellbook-darwin-arm *.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build \
		-o bin/spellbook-darwin *.go

build-win:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build \
		-o bin/spellbook-win *.go

run-darwin-arm-server: build-darwin-arm
	mv bin/spellbook-darwin-arm ./spellbook
	./spellbook server start

run-darwin-arm-cli: build-darwin-arm
	mv bin/spellbook-darwin-arm ./spellbook
