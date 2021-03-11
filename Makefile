.PHONY: clean run deploy build.local build.linux

BINARY        ?= egsam
SOURCES       = main.go
VERSION       := $(shell date '+%Y%m%d%H%M%S')
IMAGE         ?= deploy.glv.one/pitr/$(BINARY)
TAG           ?= $(VERSION)
DOCKERFILE    ?= Dockerfile
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -w -s

default: run

clean:
	rm -rf build

run: build.local
	./build/$(BINARY)

build.local: build/$(BINARY)
build.linux: build/linux/$(BINARY)

build/$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

build/linux/$(BINARY): $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(BINARY) -ldflags "$(LDFLAGS)" .

build.docker: build.linux
	docker build --rm -t "$(IMAGE):$(TAG)" -f $(DOCKERFILE) .

deploy: build.docker
	docker push "$(IMAGE):$(TAG)"

run13: tls13/main13.go
	CGO_ENABLED=0 go build -o build/$(BINARY)13 $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" tls13/main13.go
	./build/$(BINARY)13

deploy13: tls13/main13.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(BINARY)13 -ldflags "$(LDFLAGS)" tls13/main13.go
	docker build --rm -t "$(IMAGE)13:$(TAG)" -f tls13/$(DOCKERFILE) .
	docker push "$(IMAGE)13:$(TAG)"
