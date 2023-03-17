# Go Zones Makefile

This Makefile contains targets to build, run and manage Go Zones instances using Podman container engine.

## Usage
* `build-image-full:` Builds a Go Zones container image with full DNS support
* `file-to-bind-build-image:` Builds a Go Zones container image with a specified configuration file
* `server-build-image:` Builds a Go Zones container image with a specified configuration file
* `copy-config-dir-full:` Copies the default configuration files for a full DNS instance
* `copy-file-to-bind:` Copies the default configuration file for a file to BIND instance
* `copy-config-dir-server:` Copies the default configuration file for a server instance
* `start-server-instance:` Starts a Go Zones server instance
* `start-file-to-bind:` Starts a Go Zones file to BIND instance
* `start-full-instance:` Starts a Go Zones instance with full DNS support 
* `clean:` Removes unused container images and deletes the built Go Zones image
* `stop:` Stops and removes all running Go Zones instances
This Makefile contains targets to build, run and manage Go Zones instances using Podman container engine.

## Variables
* `APP_VERSION:` The version of the Go Zones application
* `GO_BUILD_CMD:` The Go build command to use
* `BINARY_NAME:` The name of the built Go Zones binary

## Examples
**Build a Go Zones image with full DNS support**
```bash
make build-image-full
```

**Build a Go Zones image with a specified configuration file**
```bash
make file-to-bind-build-image
```

**Build a Go Zones image with a specified configuration file**
```bash
make server-build-image
```
**Copy the default configuration files for a full DNS instance**
```bash
make copy-config-dir-full
```
**Copy the default configuration file for a file to BIND instance**
```bash
make copy-file-to-bind
```

**Copy the default configuration file for a server instance**
```bash
make copy-config-dir-server
```

**Start a Go Zones server instance**
```bash
make start-server-instance
```

**Start a Go Zones file to BIND instance**
```bash
make start-file-to-bind
```

**Start a Go Zones instance with full DNS support**
```bash
make start-full-instance
```

**Remove unused container images and delete the built Go Zones image**
```bash
make clean
```

**Stop and remove all running Go Zones instances**
```bash
make stop
```