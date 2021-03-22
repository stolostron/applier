# Copyright Contributors to the Open Cluster Management project

SCRIPTS_PATH ?= build

# Install software dependencies
INSTALL_DEPENDENCIES ?= ${SCRIPTS_PATH}/install-dependencies.sh
# The command to run to execute unit tests
UNIT_TEST_COMMAND ?= ${SCRIPTS_PATH}/run-unit-tests.sh

BEFORE_SCRIPT := $(shell build/before-make.sh)

export PROJECT_DIR            = $(shell 'pwd')
export PROJECT_NAME			  = $(shell basename ${PROJECT_DIR})
	
export GOPACKAGES ?= ./pkg/...
export KUBEBUILDER_HOME := /usr/local/kubebuilder

export PATH := ${PATH}:${KUBEBUILDER_HOME}/bin

.PHONY: deps
deps:
	$(INSTALL_DEPENDENCIES)

.PHONY: check
check: check-copyright

.PHONY: check-copyright
check-copyright:
	@build/check-copyright.sh

.PHONY: test
## Runs go unit tests
test:
	$(UNIT_TEST_COMMAND);

.PHONY: go/gosec-install
## Installs latest release of Gosec
go/gosec-install:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOPATH)/bin


.PHONY: go-bindata
go-bindata:
	@if which go-bindata > /dev/null; then \
		echo "##### Updating go-bindata..."; \
		cd $(mktemp -d) && GOSUMDB=off go get -u github.com/go-bindata/go-bindata/...; \
	fi
	@go-bindata --version
	go-bindata -nometadata -pkg bindata -o examples/applier/bindata/bindata_generated.go -prefix examples/applier/resources/yamlfilereader  examples/applier/resources/yamlfilereader/...

.PHONY: examples
examples:
	@mkdir -p examples/bin
	go build -o examples/bin/apply-some-yaml examples/applier/apply-some-yaml/main.go
	go build -o examples/bin/apply-yaml-in-dir examples/applier/apply-yaml-in-dir/main.go
	go build -o examples/bin/render-list-yaml examples/applier/render-list-yaml/main.go
	go build -o examples/bin/render-yaml-in-dir examples/applier/render-yaml-in-dir/main.go
	
.PHONY: build
build: 
	go install ./cmd/applier

.PHONY: install
install: build

.PHONY: oc-plugin
oc-plugin: build
	mv ${GOPATH}/bin/cm ${GOPATH}/bin/oc_cm

.PHONY: kubectl-plugin
kubectl-plugin: build
	mv ${GOPATH}/bin/cm ${GOPATH}/bin/kubectl_cm

.PHONY: functional-test
functional-test:
	ginkgo -tags functional -v --slowSpecThreshold=30 test/functional -- -v=1

.PHONY: functional-test-full
functional-test-full: 
	@build/run-functional-tests.sh

.PHONY: kind-cluster-setup
kind-cluster-setup: 
	@echo "No setup to do"
