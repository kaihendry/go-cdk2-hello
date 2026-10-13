VERSION ?= $(shell git describe --always --dirty --abbrev=7)

.PHONY: build clean
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w -X main.version=$(VERSION)" -o bootstrap .
	zip -j -X function.zip bootstrap

clean:
	rm -f bootstrap function.zip
