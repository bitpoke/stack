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
	make -C git-webhook lint

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
    # install ginkgo
	GOBIN=$(BINDIR) go install ./vendor/github.com/onsi/ginkgo/ginkgo
	# install golangci-lint
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- -b $(BINDIR) v1.16.0
	# install yq
	curl -sfL https://github.com/mikefarah/yq/releases/download/2.1.1/yq_$(OS)_$(ARCH) -o $(BINDIR)/yq
	chmod +x $(BINDIR)/yq
