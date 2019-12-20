#!/bin/bash

set -xe

export GO111MODULE="off"
export FUZZIT_API_KEY=b35543798ebb0ec63ecf3862b46691e1a2c760ba1157999bb9e8d43a22cbdd5bfeb2162dcb316261a352a6befac825d4

go get -t -v ./...

## Install go-fuzz
go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build

wget -q -O fuzzitbin https://github.com/fuzzitdev/fuzzit/releases/download/v2.4.52/fuzzit_Linux_x86_64
chmod a+x fuzzitbin

go-fuzz-build -libfuzzer -o mph.a ./pkg/mph
clang -fsanitize=fuzzer mph.a -o mph

## upload fuzz target for long fuzz testing on fuzzit.dev server or run locally for regression
./fuzzitbin create job --type $1 fuzzitdev/mph mph
