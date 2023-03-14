APP_VERSION=$(shell hack/version.sh)
GO_BUILD_CMD= CGO_ENABLED=0 go build -ldflags="-X main.appVersion=$(APP_VERSION)"

BINARY_NAME=go-zones

.PHONY: build-image
build-image:
	podman build -f Containerfile -t ${BINARY_NAME} .