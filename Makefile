SHELL := /bin/bash -o pipefail
PROJDIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
PROJNAME := $(shell basename $(PROJDIR))
VERSION := $(shell git describe --abbrev=0 --tags)

versionDir=github.com/prometheus/common/version
gitBranch=$(shell git rev-parse --abbrev-ref HEAD)
gitCommit=$(shell git rev-parse --short HEAD)
gitUser=$(shell git config user.name | tr -d ' ')
buildDate=$(shell TZ=Asia/Shanghai date +%FT%T%z)

ldflags="-s -w -X ${versionDir}.Version=${VERSION} -X ${versionDir}.Revision=${gitCommit} \
    -X ${versionDir}.Branch=${gitBranch} -X ${versionDir}.BuildUser=${gitUser} \
    -X ${versionDir}.BuildDate=${buildDate}"
all: binary

binary:
	@echo "build the ${PROJNAME}"
	go build -ldflags ${ldflags}
	@echo "build done."

docker: binary
	docker build . -t $(PROJNAME):$(gitCommit)

.PHONY: doc-build
doc-build:
	go doc utils
	go doc collector
	go doc collector_s
