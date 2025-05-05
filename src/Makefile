.PHONY: clean

function.zip: main.go
	# Build the Go binary, stripping debug info and symbols for smaller size
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bootstrap
	# Zip only the executable
	@zip -j function.zip bootstrap

clean:
	@rm -f bootstrap function.zip
