.PHONY: build
build:
	docker-compose -f docker/docker-compose.dev.yml build

.PHONY: tag
tag:
	docker tag cloudflare-spectrum-ddns:dev connormckelvey/cloudflare-spectrum-ddns:$$DOCKER_TAG

.PHONY: push
push:
	docker push connormckelvey/cloudflare-spectrum-ddns:$$DOCKER_TAG

.PHONY: release
release: build tag push

.PHONY: latest
latest:  
	DOCKER_TAG=latest make release

.PHONY: test
test: mocks
	go test -count=1 -v ./...

.PHONY: mocks
mocks:
	mockgen -source=pkg/ipchecker/dnsclient.go -package mocks -destination pkg/ipchecker/mocks/dnsclint.go