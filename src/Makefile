build:
	@GOOS=linux GOARCH=arm64 go build -o bootstrap
	@zip -j function.zip bootstrap

clean:
	@rm -f bootstrap function.zip
