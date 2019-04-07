## Makefile

build-suggest:
	go build -o build/suggest ./pkg/cmd/suggest/

build: build-suggest

build-docker:
	docker build --no-cache -t suggest:latest .

test:
	go test ./pkg/...

clean:
	rm -rf build
