/*
Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package provider contains the cloud provider specific implementations to manage machines
package provider

import (
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/metal-stack/machine-controller-manager-provider-metal/pkg/spi"
	metalgo "github.com/metal-stack/metal-go"
)

// Provider is the struct that implements the driver interface
// It is used to implement the basic driver functionalities
type Provider struct {
	SPI spi.SessionProviderInterface
}

// NewProvider returns an empty provider object
func NewProvider(spi spi.SessionProviderInterface) driver.Driver {
	return &Provider{
		SPI: spi,
	}
}

func (p *Provider) initDriver(secret *corev1.Secret) (*metalgo.Driver, error) {
	token := strings.TrimSpace(string(secret.Data["metalAPIKey"]))
	hmac := strings.TrimSpace(string(secret.Data["metalAPIHMac"]))

	u := secret.Data["metalAPIURL"]
	url := strings.TrimSpace(string(u))

	return metalgo.NewDriver(url, token, hmac)
}
