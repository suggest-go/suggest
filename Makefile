## Makefile

.PHONY: build test vet clean

BUILD_FLAGS = $(GO_BUILD_FLAGS)

default: build

build-suggest:
	go build $(BUILD_FLAGS) -o build/suggest ./cmd/suggest/

build-lm:
	go build $(BUILD_FLAGS) -o build/lm ./cmd/language-model/

build-spellchecker:
	go build $(BUILD_FLAGS) -o build/spellchecker ./cmd/spellchecker/

build: build-suggest build-lm build-spellchecker

build-docker:
	docker build --no-cache -t suggest:latest .

test:
	go test -race -v ./...

vet:
	go vet ./...

clean:
	rm -rf build
