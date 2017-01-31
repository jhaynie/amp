
.PHONY: install

install:
	@go install ./cmd/demo-publisher
	@go install ./cmd/demo-subscriber

image:
	@docker build -f Dockerfile.publisher -t appcelerator/amp-demo-publisher .
	@docker build -f Dockerfile.subscriber -t appcelerator/amp-demo-subscriber .
