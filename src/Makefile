VERSION ?= $(shell git describe --always --dirty --abbrev=7)
OUT_DIR ?= .

.PHONY: all bundle clean
all: bundle
	zip -j function.zip bootstrap

# CDK bundles directly into its asset directory; no Docker or shared bucket needed.
bundle:
	mkdir -p "$(OUT_DIR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w -X main.version=$(VERSION)" -o "$(OUT_DIR)/bootstrap" .

clean:
	rm -f bootstrap function.zip
