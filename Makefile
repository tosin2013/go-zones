APP_VERSION=$(shell hack/version.sh)
GO_BUILD_CMD= CGO_ENABLED=0 go build -ldflags="-X main.appVersion=$(APP_VERSION)"

BINARY_NAME=go-zones

.PHONY: build-image-full
build-image-full:
	sudo podman build -f Containerfile-to-BIND-FULL -t ${BINARY_NAME} .

.PHONY: file-to-bind-build-image
file-to-bind-build-image:
	sudo podman build -f Containerfile-to-BIND -t ${BINARY_NAME} .

.PHONY: server-build-image
server-build-image:
	sudo podman build -f Containerfile -t ${BINARY_NAME} .

.PHONY: copy-config-dir-full
copy-config-dir-full:
	mkdir -p config && cp example.config.yml config/config.yml &&  yq -o=json example.server.yml > config/server.yml

.PHONY: copy-file-to-bind
copy-file-to-bind:
	mkdir -p config && cp example.server.yml config/server.yml

.PHONY: copy-config-dir-server
copy-config-dir-server:
	mkdir -p config && cp example.config.yml config/config.yml

.PHONY: start-server-instance
start-server-instance:
	sudo podman run --name go-zones -p 8080:8080 -p -v $(CURDIR)/config:/etc/go-zones/:Z localhost/${BINARY_NAME}

.PHONY: start-file-to-bind
start-file-to-bind:
	sudo podman run --name go-zones-ftb -p 8053:8053 -p -v $(CURDIR)/config:/etc/go-zones/:Z localhost/${BINARY_NAME}

.PHONY: start-full-instance
start-full-instance:
	sudo podman run --name go-zones-full -p 8080:8080 -p 8053:53  -v $(CURDIR)/config:/etc/go-zones/:Z localhost/${BINARY_NAME}

.PHONY: clean
clean:
	sudo podman image prune --force --filter "label=io.containers.image.dangling=true"
	sudo podman rmi -f ${BINARY_NAME}

.PHONY: stop
stop:
	sudo podman stop $(shell sudo podman ps -a -q)
	sudo podman rm $(shell sudo podman ps -a -q)


