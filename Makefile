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

MKDIR = mkdir -p
RM = rm -rvf
GO = go
GLIDE = $(GOPATH)/bin/glide

sources := $(wildcard *.go)
coverage = out
coverage_out = $(coverage)/coverage.out
coverage_html = $(coverage)/coverage.html

all: prereqs fmt vet test
.PHONY: all

prereqs: $(GLIDE)
	${GLIDE} install
.PHONY: prereqs

$(GLIDE):
	${GO} get github.com/Masterminds/glide

fmt:
	${GO} fmt
.PHONY: fmt

vet:
	${GO} vet
.PHONY: vet

test: $(coverage_html)
.PHONY: test

clean:
	${RM} out
.PHONY: clean

$(coverage_out): $(sources)
	${MKDIR} $(coverage)
	${GO} test -coverprofile=$(coverage_out)

$(coverage_html): $(coverage_out)
	${GO} tool cover -func=$(coverage_out)
	${GO} tool cover -html=$(coverage_out) -o $(coverage_html)
