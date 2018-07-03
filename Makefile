# Copyright 2018 Platform 9 Systems, Inc.

# Define some constants
#######################
ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR ?= build
BIN_DIR	?= bin
SSHPROVIDER_PKG = github.com/platform9/ssh-provider
TYPES_FILES    = $(shell find pkg/apis -name types.go)

build: .generate_files

.generate_exes: $(BIN_DIR)/defaulter-gen \
                $(BIN_DIR)/deepcopy-gen \
                $(BIN_DIR)/conversion-gen \
                $(BIN_DIR)/client-gen \
                $(BIN_DIR)/lister-gen \
                $(BIN_DIR)/informer-gen \
                $(BIN_DIR)/openapi-gen
	touch $@

# Regenerate all files if the gen exes changed or any "types.go" files changed
.generate_files: .generate_exes $(TYPES_FILES)
	# generate apiserver deps
	$(BUILD_DIR)/update-apiserver-gen.sh
	# generate all pkg/client contents
	$(BUILD_DIR)/update-client-gen.sh
	touch $@

$(BIN_DIR)/defaulter-gen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/k8s.io/code-generator/cmd/defaulter-gen

$(BIN_DIR)/deepcopy-gen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/k8s.io/code-generator/cmd/deepcopy-gen

$(BIN_DIR)/conversion-gen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/k8s.io/code-generator/cmd/conversion-gen

$(BIN_DIR)/client-gen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/k8s.io/code-generator/cmd/client-gen

$(BIN_DIR)/lister-gen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/k8s.io/code-generator/cmd/lister-gen

$(BIN_DIR)/informer-gen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/k8s.io/code-generator/cmd/informer-gen

$(BIN_DIR)/openapi-gen: vendor/k8s.io/code-generator/cmd/openapi-gen
	go build -o $@ $(SSHPROVIDER_PKG)/$^

.PHONY: generate_mocks

# Regenerate all mocks if the mockgen binary, or sources have changed
generate_mocks: $(BIN_DIR)/mockgen \
				pkg/machine/mock/client_generated.go

$(BIN_DIR)/mockgen:
	go build -o $@ $(SSHPROVIDER_PKG)/vendor/github.com/golang/mock/mockgen

pkg/machine/mock/client_generated.go: pkg/machine/client.go
	mockgen -source=$^ -package=mock -destination=$@
