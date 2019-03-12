NAME=okchain/baseimage
VERSION=$(shell cat ./release)
ARCH=$(shell uname -m)
DOCKER_TAG ?= $(ARCH)-$(VERSION)

DOCKER_BASE_x86_64=ubuntu:xenial

DOCKER_BASE=$(DOCKER_BASE_$(ARCH))

ifeq ($(DOCKER_BASE), )
$(error "Architecture \"$(ARCH)\" is unsupported")
endif

all: docker

Dockerfile: Dockerfile.in Makefile
	@echo "# Generated from Dockerfile.in.  DO NOT EDIT!" > $@
	@cat Dockerfile.in | \
	sed -e  "s|_DOCKER_BASE_|$(DOCKER_BASE)|" >> $@

docker: Dockerfile release
	@echo "Generating docker"
	@docker build -t $(NAME):$(DOCKER_TAG) .

clean:
	-rm Dockerfile
