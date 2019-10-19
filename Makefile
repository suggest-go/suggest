## Makefile

.PHONY: build test vet clean

BUILD_FLAGS = -mod=vendor $(GO_BUILD_FLAGS)

default: build

build-suggest:
	go build $(BUILD_FLAGS) -o build/suggest ./pkg/cmd/suggest/

build-lm:
	go build $(BUILD_FLAGS) -o build/lm ./pkg/cmd/language-model/

build: download test vet build-suggest build-lm
build-bin: download build-suggest build-lm

build-docker:
	docker build --no-cache -t suggest:latest .

test:
	go test -race -v ./pkg/...

download:
	go mod download
	go mod vendor

vet:
	go vet ./pkg/...

clean:
	rm -rf build
