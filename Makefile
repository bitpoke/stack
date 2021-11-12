# Project Setup
PROJECT_NAME := bitpoke-stack
PROJECT_REPO := github.com/bitpoke/stack

PLATFORMS := linux_amd64 darwin_amd64

include build/makelib/common.mk

GO111MODULE = on
GO_PROJECT := $(PROJECT_REPO)
GO_SUBDIRS := default-backend git-webhook
GO_STATIC_PACKAGES := $(GO_PROJECT)/git-webhook $(GO_PROJECT)/default-backend
include build/makelib/golang.mk
include build/makelib/kubebuilder-v2.mk

IMAGES := stack-default-backend stack-git-webhook
DOCKER_REGISTRY ?= docker.io/bitpoke
include build/makelib/image.mk

HELM_CHARTS := stack wordpress-site
include build/makelib/helm.mk
.PHONY: .helm.publish
.helm.publish:
	@$(INFO) publishing helm charts
	@rm -rf $(WORK_DIR)/charts
	@git clone -q git@github.com:bitpoke/helm-charts.git $(WORK_DIR)/charts
	@cp $(HELM_OUTPUT_DIR)/*.tgz $(WORK_DIR)/charts/docs/
	@git -C $(WORK_DIR)/charts add $(WORK_DIR)/charts/docs/*.tgz
	@git -C $(WORK_DIR)/charts commit -q -m "Added $(call list-join,$(COMMA)$(SPACE),$(foreach c,$(HELM_CHARTS),$(c)-v$(HELM_CHART_VERSION)))"
	@git -C $(WORK_DIR)/charts push -q
	@$(OK) publishing helm charts
.publish.run: .helm.publish

.PHONY: .helm.add-repos
.helm.add-repos: $(HELM)
	@$(HELM) repo add bitpoke https://helm-charts.bitpoke.io
	@$(HELM) repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
.helm.lint.stack: .helm.add-repos

#
# custom targets
#
.PHONY: .collect-crds
.collect-crds: helm.dep |$(HELM)
	@$(INFO) collecting CRDs from dependent charts
	@rm -rf $(WORK_DIR)/collect-crds deploy/00-crds.yaml
	@mkdir -p $(WORK_DIR)/collect-crds
	@$(HELM) template --include-crds --output-dir $(WORK_DIR)/collect-crds deploy/charts/stack >/dev/null
	@for crd in $(WORK_DIR)/collect-crds/stack/charts/*/crds/* ; do \
		echo "---" >> deploy/00-crds.yaml; \
		cat $${crd} >> deploy/00-crds.yaml; \
	done
	@$(OK) collecting CRDs from dependent charts
.generate.run: .collect-crds
