APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)
BINDIR ?= $(PWD)/bin
CHARTDIR ?= $(PWD)/charts

OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= amd64

PATH := $(BINDIR):$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

.PHONY: charts
charts:
	yq w -i $(CHARTDIR)/stack/Chart.yaml version "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/stack/Chart.yaml appVersion "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/stack/values.yaml nginx-ingress.defaultBackend.image.tag "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/wordpress-site/Chart.yaml version "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/wordpress-site/Chart.yaml appVersion "$(APP_VERSION)"

lint:
	helm lint charts/stack
	helm lint charts/wordpress-site --set 'site.domains[0]=example.com'
	helm dep build charts/wordpress-site
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
