MAKEFLAGS += --warn-undefined-variables
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: all prereq fmt vet test clean

sources := $(wildcard *.go)

all: prereq fmt vet test

prereqs:
	glide install

fmt:
	go fmt

vet:
	go vet

test: $(sources) out/coverage.html

out:
	mkdir -p out

clean:
	rm -rfv out

out/coverage.out: out
	go test -coverprofile=out/coverage.out

out/coverage.html: out/coverage.out
	go tool cover -func=out/coverage.out
	go tool cover -html=out/coverage.out -o out/coverage.html
