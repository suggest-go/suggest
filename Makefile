## Makefile

build-indexer:
	go build -o build/indexer ./pkg/cmd/indexer/

build-suggest:
	go build -o build/suggest ./pkg/cmd/suggest/

build: build-indexer build-suggest

build-docker:
	docker build --no-cache -t suggest:0.0.1 .

test:
	go test ./pkg/...

clean:
	rm -rf build
