SHELL = /bin/bash -o pipefail
PROJDIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
PROJNAME := $(shell basename $(PROJDIR))
GITHASH := $(shell git rev-parse --short HEAD)

all: binary

binary:
	go build

docker: binary
	docker build . -t $(PROJNAME):$(GITHASH)
