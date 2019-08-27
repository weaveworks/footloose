UID_GID?=$(shell id -u):$(shell id -g)
GO_VERSION="1.12.6"

all: binary

binary: vendor
	docker run -it --rm -v $(shell pwd):/build -w /build golang:${GO_VERSION} sh -c "\
		make footloose && \
		chown -R ${UID_GID} bin"

footloose: bin/footloose
bin/footloose:
	CGO_ENABLED=0 go build -mod=vendor -o bin/footloose .

D := $(shell go env GOPATH)/bin
install: bin/footloose
	mkdir -p $(D)
	cp $^ $(D)

vendor:
	go mod vendor

.PHONY: bin/footloose install binary footloose vendor

