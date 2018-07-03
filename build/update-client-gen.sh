#!/bin/bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The only argument this script should ever be called with is '--verify-only'

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

REPO_ROOT=$(cd $(dirname "${BASH_SOURCE}")/..; pwd)
BINDIR=${REPO_ROOT}/bin

# Generate the internal clientset (pkg/client/clientset_generated/internalclientset)
${BINDIR}/client-gen "$@" \
		--input-base "github.com/platform9/ssh-provider/pkg/apis/" \
		--input sshprovider/ \
		--clientset-path "github.com/platform9/ssh-provider/pkg/client/clientset_generated/" \
		--clientset-name internalclientset \
		--go-header-file "boilerplate.go.txt"
# Generate the versioned clientset (pkg/client/clientset_generated/clientset)
${BINDIR}/client-gen "$@" \
		--input-base "github.com/platform9/ssh-provider/pkg/apis/" \
		--input "sshprovider/v1alpha1" \
		--clientset-path "github.com/platform9/ssh-provider/pkg/client/clientset_generated/" \
		--clientset-name "clientset" \
		--go-header-file "boilerplate.go.txt"
# generate lister
	${BINDIR}/lister-gen "$@" \
		--input-dirs="github.com/platform9/ssh-provider/pkg/apis/sshprovider" \
		--input-dirs="github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1" \
		--output-package "github.com/platform9/ssh-provider/pkg/client/listers_generated" \
		--go-header-file "boilerplate.go.txt"
# generate informer
${BINDIR}/informer-gen "$@" \
		--input-dirs "github.com/platform9/ssh-provider/pkg/apis/sshprovider" \
		--input-dirs "github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1" \
		--internal-clientset-package "github.com/platform9/ssh-provider/pkg/client/clientset_generated/internalclientset" \
		--versioned-clientset-package "github.com/platform9/ssh-provider/pkg/client/clientset_generated/clientset" \
		--listers-package "github.com/platform9/ssh-provider/pkg/client/listers_generated" \
		--output-package "github.com/platform9/ssh-provider/pkg/client/informers_generated" \
		--go-header-file "boilerplate.go.txt"