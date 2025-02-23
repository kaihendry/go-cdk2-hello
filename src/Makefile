build:
	@GOOS=linux GOARCH=amd64 go build -o bootstrap
	@zip -j function.zip bootstrap

clean:
	@rm -f bootstrap function.zip