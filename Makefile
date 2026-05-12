BINARY := wd-go
VERSION := $(shell grep 'const version' cmd/version.go | sed 's/.*"\(.*\)".*/\1/')
LDFLAGS_PROD := -ldflags="-s -w -X wd-go/cmd.version=$(VERSION)"
GCFLAGS_DEV  := -gcflags="all=-N -l"

.PHONY: dev prod clean test install

dev:
	go build $(GCFLAGS_DEV) -o $(BINARY) .

prod:
	go build $(LDFLAGS_PROD) -o $(BINARY) .

clean:
	rm -f $(BINARY)

test:
	go test ./...

install: prod
	go install $(LDFLAGS_PROD) .
