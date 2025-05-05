test:
	curl https://stghello.dabase.com

.PHONY: mb
mb:
	-aws s3 mb s3://hendry-lambdas/

deploy:
	cdk deploy --outputs-file outputs.json