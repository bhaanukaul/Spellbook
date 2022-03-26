

all: test build-cli build-server

test:

build-server: build-linux-server build-darwin-arm-server build-darwin-server build-win-server

build-cli: build-linux-cli build-darwin-arm-cli build-darwin-cli build-win-cli

build-linux-cli:
	GOOS=linux GOARCH=amd64 go build \
		-o bin/cli/Spellbook-linux cmd/cli/*.go

build-darwin-arm-cli:
	GOOS=darwin GOARCH=arm64 go build \
		-o bin/cli/Spellbook-darwin-arm cmd/cli/*.go

build-darwin-cli:
	GOOS=darwin GOARCH=amd64 go build \
		-o bin/cli/Spellbook-darwin cmd/cli/*.go

build-win-cli:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build \
		-o bin/cli/Spellbook-win cmd/cli/*.go

build-linux-server:
	GOOS=linux GOARCH=amd64 go build \
		-o bin/server/Spellbook-Server-linux cmd/server/*.go

build-darwin-arm-server:
	GOOS=darwin GOARCH=arm64 go build \
		-o bin/server/Spellbook-Server-darwin-arm cmd/server/*.go

build-darwin-server:
	GOOS=darwin GOARCH=amd64 go build \
		-o bin/server/Spellbook-Server-darwin cmd/server/*.go

build-win-server:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build \
		-o bin/server/Spellbook-Server-win cmd/server/*.go