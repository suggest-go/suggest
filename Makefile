## Makefile

build-suggest:
	go build -o build/suggest ./pkg/cmd/suggest/

build-lm:
	go build -o build/lm ./pkg/cmd/language-model/

build: build-suggest build-lm

build-docker:
	docker build --no-cache -t suggest:latest .

test:
	go test ./pkg/...

clean:
	rm -rf build
