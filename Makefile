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
	yq w -i $(CHARTDIR)/stack/values.yaml nginx-ingress.defaultBackend.image.tag "$(APP_VERSION:v%=%)"
	yq w -i $(CHARTDIR)/stack/values.yaml git-webhook.image.tag "$(APP_VERSION:v%=%)"
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

test:
	make -C git-webhook test


define getVersion
$(shell python -c "import yaml; print([x['version'] for x in yaml.load(open('charts/stack/requirements.lock', 'r'))['dependencies'] if x['name'] == '$1'  ][0])")
endef

MANIFESTS_DIR ?= deploy/manifests
CRDS_FILE ?= $(MANIFESTS_DIR)/00-crds.yaml

CERT_MANAGER_TAG ?= $(call getVersion,cert-manager)
MYSQL_OPERATOR_TAG ?= v$(call getVersion,mysql-operator)
WORDPRESS_OPERATOR_TAG ?= $(call getVersion,wordpress-operator)

.PHONY: collect-crds
collect-crds:
	$(info ---- CERT_MANAGER_TAG = $(CERT_MANAGER_TAG))
	$(info ---- WORDPRESS_OPERATOR_TAG = $(WORDPRESS_OPERATOR_TAG))
	$(info ---- MYSQL_OPERATOR_TAG = $(MYSQL_OPERATOR_TAG))

	# wordpress operator
	wget https://raw.githubusercontent.com/presslabs/wordpress-operator/$(WORDPRESS_OPERATOR_TAG)/config/crds/wordpress_v1alpha1_wordpress.yaml -O - > $(CRDS_FILE)

	# mysql operator
	wget https://raw.githubusercontent.com/presslabs/mysql-operator/$(MYSQL_OPERATOR_TAG)/config/crds/mysql_v1alpha1_mysqlcluster.yaml -O - >>  $(CRDS_FILE)
	wget https://raw.githubusercontent.com/presslabs/mysql-operator/$(MYSQL_OPERATOR_TAG)/config/crds/mysql_v1alpha1_mysqlbackup.yaml -O  - >> $(CRDS_FILE)

	# cert manager
	wget https://raw.githubusercontent.com/jetstack/cert-manager/$(CERT_MANAGER_TAG)/deploy/manifests/00-crds.yaml -O - >> $(CRDS_FILE)
