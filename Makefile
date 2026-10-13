.PHONY: test synth deploy clean

test:
	cd src && go test -race ./...
	cd src && go vet ./...
	npm run check
	npm test
	bash -n deploy.sh scripts/smoke-test.sh

synth:
	npm run synth

deploy:
	npm run deploy -- --require-approval never
	./scripts/smoke-test.sh

clean:
	$(MAKE) -C src clean
	rm -rf cdk.out
