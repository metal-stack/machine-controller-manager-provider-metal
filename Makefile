# Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BINARY_PATH         := bin/
COVERPROFILE        := test/output/coverprofile.out
IMAGE_REPOSITORY    := ghcr.io/metal-stack/machine-controller-manager-provider-metal
IMAGE_TAG           := $(or ${GITHUB_TAG_NAME}, latest)
PROVIDER_NAME       := MetalProvider
PROJECT_NAME        := gardener
CONTROL_NAMESPACE  := default
CONTROL_KUBECONFIG := dev/target-kubeconfig.yaml
TARGET_KUBECONFIG  := dev/target-kubeconfig.yaml
VERSION            := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))

#########################################
# Rules for testing
#########################################

.PHONY: test-unit
test-unit:
	.ci/test

#########################################
# Rules for build/release
#########################################

.PHONY: release
release: build-local build docker-image docker-login docker-push rename-binaries

.PHONY: build-local
build-local:
	@env LOCAL_BUILD=1 .ci/build

.PHONY: build
build:
	@GO111MODULE=on go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -o bin/machine-controller \
    -ldflags "-X main.version=${VERSION}" \
    cmd/machine-controller/main.go
	strip bin/machine-controller

.PHONY: docker-image
docker-image:
	@docker build -t $(IMAGE_REPOSITORY):$(IMAGE_TAG) .

.PHONY: docker-push
docker-push:
	@docker push $(IMAGE_REPOSITORY):$(IMAGE_TAG)

.PHONY: rename-binaries
rename-binaries:
	@if [[ -f bin/machine-controller ]]; then cp bin/machine-controller machine-controller-darwin-amd64; fi
	@if [[ -f bin/rel/machine-controller ]]; then cp bin/rel/machine-controller machine-controller-linux-amd64; fi

.PHONY: clean
clean:
	@rm -rf bin/
	@rm -f *linux-amd64
	@rm -f *darwin-amd64
