GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
BUILD_DIR = dist/${GOOS}_${GOARCH}

ifeq ($(GOOS),windows)
OUTPUT_PATH = ${BUILD_DIR}/baton-cloudamqp.exe
else
OUTPUT_PATH = ${BUILD_DIR}/baton-cloudamqp
endif

.PHONY: build
build:
	rm -f ${OUTPUT_PATH}
	mkdir -p ${BUILD_DIR}
	go build -o ${OUTPUT_PATH} cmd/baton-cloudamqp/*.go

.PHONY: update-deps
update-deps:
	go get -d -u ./...
	go mod tidy -v
	go mod vendor

.PHONY: add-dep
add-dep:
	go mod tidy -v
	go mod vendor

.PHONY: lint
lint:
	golangci-lint run
