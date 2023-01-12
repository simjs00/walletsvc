NAMESPACE = `echo walletsvc`
BUILD_TIME = `date +%FT%T%z`
BUILD_VERSION = `git describe --tag --abbrev=0`
COMMIT_HASH = `git rev-parse --short HEAD`

.PHONY: build
build:
	@go build -v -tags dynamic -ldflags "-X main.Namespace=${NAMESPACE} \
	  -X main.BuildTime=${BUILD_TIME} \
	  -X main.BuildVersion=${BUILD_VERSION} \
	  -X main.CommitHash=${COMMIT_HASH}" \
	  -race -o ./build/app ./main

.PHONY: kill-process
kill-process:
	@lsof -i :8099 | awk '$$1 ~ /app/ { print $$2 }' | xargs kill -9 || true

.PHONY: run
run: kill-process build
	@./build/app