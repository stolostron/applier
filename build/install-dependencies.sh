#!/bin/bash -e

# Copyright Contributors to the Open Cluster Management project

# Go tools
_OS=$(go env GOOS)
_ARCH=$(go env GOARCH)
KubeBuilderVersion="2.2.0"

if ! which patter > /dev/null; then      echo "Installing patter ..."; GO111MODULE=off go get -u github.com/apg/patter; fi
if ! which gocovmerge > /dev/null; then  echo "Installing gocovmerge..."; GO111MODULE=off go get -u github.com/wadey/gocovmerge; fi
if ! which golangci-lint > /dev/null; then
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.6
fi

if ! which kubebuilder > /dev/null; then
   # Install kubebuilder for unit test
   echo "Install Kubebuilder components for test framework usage!"

   # download kubebuilder and extract it to tmp
   curl -L https://go.kubebuilder.io/dl/"$KubeBuilderVersion"/"${_OS}"/"${_ARCH}" | tar -xz -C /tmp/

   # move to a long-term location and put it on your path
   # (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
   sudo mv /tmp/kubebuilder_"$KubeBuilderVersion"_"${_OS}"_"${_ARCH}" $KUBEBUILDER_HOME
fi

# Build tools

# Image tools

# Check tools
