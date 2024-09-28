EXECUTABLE=gfping
# COMMIT=$(shell git rev-parse HEAD)
# VERSION=$(shell git describe --tags --exact-match --always)
# VERSION=$(shell git describe --tags --always --long --dirty)
VERSION=0.1
DATE=$(shell date +'%F-%H%M%S')
LINUX=$(EXECUTABLE)_linux_amd64
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
DARWIN=$(EXECUTABLE)_darwin_amd64

STATIC:=0

ifeq ($(STATIC),1)
LDFLAGS+=-s -w -extldflags "-static"
endif

all: linux windows darwin

build:
	@echo version: $(VERSION)
	go build -o $(EXECUTABLE)


linux: $(LINUX)

windows: $(WINDOWS)

darwin: $(DARWIN)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -o ./bin/$(WINDOWS) -ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.Date=$(DATE)"

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -o ./bin/$(LINUX) -ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.Date=$(DATE)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/$(DARWIN) -ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.Date=$(DATE)"


# Cleans our projects: deletes binaries
clean:
	go clean
	rm -rf bin

.PHONY: all clean build
