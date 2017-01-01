# Makefile
#
# Author::    Chirantan Mitra
# Copyright:: Copyright (c) 2015-2017. All rights reserved
# License::   MIT

MAKEFLAGS += --warn-undefined-variables
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

ifndef GOPATH
$(error GOPATH not set)
endif

sources := $(wildcard *.go)

all: prereq fmt vet test
.PHONY: all

prereqs:
	glide install
.PHONY: prereq

fmt:
	go fmt
.PHONY: fmt

vet:
	go vet
.PHONY: vet

test: out/coverage.html
.PHONY: test

out:
	mkdir -p out

clean:
	rm -rfv out
.PHONY: clean

out/coverage.out: $(sources) out
	go test -coverprofile=out/coverage.out

out/coverage.html: $(sources) out/coverage.out
	go tool cover -func=out/coverage.out
	go tool cover -html=out/coverage.out -o out/coverage.html
