## Makefile

.PHONY: build

GO_BUILD = go build -mod=vendor

default: build

build-suggest:
	$(GO_BUILD) -o build/suggest ./pkg/cmd/suggest/

build-lm:
	$(GO_BUILD) -o build/lm ./pkg/cmd/language-model/

build: download test build-suggest build-lm

build-docker:
	docker build --no-cache -t suggest:latest .

test:
	go test -race -v ./pkg/...
	go vet ./pkg/...

download:
	go mod download
	go mod vendor

clean:
	rm -rf build
