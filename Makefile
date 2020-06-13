.PHONY: clean run deploy build.local build.linux

BINARY        ?= egsam
SOURCES       = $(shell find . -name '*.go')
BUILD_FLAGS   ?= -v
PORT          ?= 1965
LDFLAGS       ?= -w -s -X main.port=$(PORT)

default: run

clean:
	rm -rf build

run: build.local
	./build/$(BINARY)

deploy: build.linux
	scp build/linux/$(BINARY) ec2-user@$(PRODUCTION):$(BINARY)-next
	ssh ec2-user@$(PRODUCTION) 'cp $(BINARY) $(BINARY)-old'
	ssh ec2-user@$(PRODUCTION) 'mv $(BINARY)-next $(BINARY)'
	ssh ec2-user@$(PRODUCTION) 'sudo systemctl restart $(BINARY)'
	scp -r static ec2-user@$(PRODUCTION):static

rollback:
	ssh ec2-user@$(PRODUCTION) 'mv $(BINARY)-old $(BINARY)'
	ssh ec2-user@$(PRODUCTION) 'sudo systemctl restart $(BINARY)'

build.local: build/$(BINARY)
build.linux: build/linux/$(BINARY)

build/$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

build/linux/$(BINARY): $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(BINARY) -ldflags "$(LDFLAGS)" .
