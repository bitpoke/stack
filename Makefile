APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)
BINDIR ?= $(PWD)/bin
CHARTDIR ?= $(PWD)/chart/stack

OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= amd64

PATH := $(BINDIR):$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

.PHONY: chart
chart:
	yq w -i $(CHARTDIR)/Chart.yaml version "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/Chart.yaml appVersion "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/values.yaml nginx-ingress.defaultBackend.image.tag "$(APP_VERSION)"

lint:
	helm lint chart/stack

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	curl -sfL https://github.com/mikefarah/yq/releases/download/2.1.1/yq_$(OS)_$(ARCH) -o $(BINDIR)/yq
	chmod +x $(BINDIR)/yq
