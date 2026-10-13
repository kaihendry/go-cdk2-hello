.PHONY: test synth deploy clean

test:
	cd src && go test -race ./... && go vet ./...
	pnpm check && pnpm test
	test -z "$$(gofmt -l src)"
	bash -n deploy.sh

synth:
	$(MAKE) -C src
	pnpm cdk synth

deploy:
	$(MAKE) -C src
	pnpm cdk deploy --require-approval never --outputs-file outputs.json

clean:
	$(MAKE) -C src clean
	rm -rf cdk.out
