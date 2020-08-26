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
IMAGE_REPOSITORY    := metal-stack/machine-controller-manager-metal
IMAGE_TAG           := $(or ${GITHUB_TAG_NAME}, latest)
PROVIDER_NAME       := MetalProvider
PROJECT_NAME        := gardener
CONTROL_NAMESPACE  := default
CONTROL_KUBECONFIG := dev/target-kubeconfig.yaml
TARGET_KUBECONFIG  := dev/target-kubeconfig.yaml

#########################################
# Rules for starting machine-controller locally
#########################################

.PHONY: start
start:
	go run \
			cmd/machine-controller/main.go \
			--control-kubeconfig=$(CONTROL_KUBECONFIG) \
			--target-kubeconfig=$(TARGET_KUBECONFIG) \
			--namespace=$(CONTROL_NAMESPACE) \
			--machine-creation-timeout=20m \
			--machine-drain-timeout=5m \
			--machine-health-timeout=10m \
			--machine-pv-detach-timeout=2m \
			--machine-safety-apiserver-statuscheck-timeout=30s \
			--machine-safety-apiserver-statuscheck-period=1m \
			--machine-safety-orphan-vms-period=30m \
			--v=3

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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -o bin/machine-controller \
    -ldflags "-X main.version=${VERSION}-$(git rev-parse HEAD)" \
    cmd/machine-controller/main.go

.PHONY: docker-image
docker-image:
	@docker build -t $(IMAGE_REPOSITORY):$(IMAGE_TAG) .

.PHONY: docker-push
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
