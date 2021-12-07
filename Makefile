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