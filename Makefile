APP_VERSION=$(shell hack/version.sh)
GO_BUILD_CMD= CGO_ENABLED=0 go build -ldflags="-X main.appVersion=$(APP_VERSION)"

BINARY_NAME=go-zones

.PHONY: build-image
build-image:
	sudo podman build -f Containerfile-to-BIND-FULL -t ${BINARY_NAME} .

.PHONY: copy-config-dir-server
copy-config-dir-server:
	mkdir -p config && cp example.config.yml config/config.yml &&  yq -o=json example.server.yml > config/server.yml

.PHONY: start-server-instance
start-server-instance:
	sudo podman run -p 8080:8080 -p -v $(CURDIR)/config:/etc/go-zones/:Z localhost/${BINARY_NAME}

.PHONY: start-full-instance
start-full-instance:
	sudo podman run -p 8080:8080 -p 53  -v $(CURDIR)/config:/etc/go-zones/:Z localhost/${BINARY_NAME}

.PHONY: clean
clean:
    -sudo podman stop  $$(sudo podman ps -aq) 
	-sudo podman rm $$(sudo podman ps -aq) 
	-sudo podman image prune --force --filter "label=io.containers.image.dangling=true"
	-sudo podman rmi -f ${BINARY_NAME}




