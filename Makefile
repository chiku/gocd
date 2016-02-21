MAKEFLAGS += --warn-undefined-variables
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: all fmt vet test

all: fmt vet test

fmt:
	go fmt

vet:
	go vet

test:
	go test
