APP_VERSION=$(shell hack/version.sh)
GO_BUILD_CMD= CGO_ENABLED=0 go build -ldflags="-X main.appVersion=$(APP_VERSION)"

BINARY_NAME=go-zones

.PHONY: build-image
build-image:
	podman build -f Containerfile -t ${BINARY_NAME} .

.PHONY: copy-config-dir
copy-config-dir:
	mkdir -p config && cp example.config.yml config/config.yml

.PHONY: start-instance
start-instance:
	podman run -p 8080:8080 -v $(CURDIR)/config:/etc/go-zones/ ${BINARY_NAME}