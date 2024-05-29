/*
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved.

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
	"context"
	"testing"

	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
)

func TestProvider_GetVolumeIDs(t *testing.T) {
	tests := []struct {
		name    string
		req     *driver.GetVolumeIDsRequest
		want    *driver.GetVolumeIDsResponse
		wantErr error
	}{
		{
			name: "valid lightbits volume",
			req: &driver.GetVolumeIDsRequest{
				PVSpecs: []*corev1.PersistentVolumeSpec{
					{

						PersistentVolumeSource: corev1.PersistentVolumeSource{
							CSI: &corev1.CSIPersistentVolumeSource{
								Driver:       "csi.lightbitslabs.com",
								VolumeHandle: "mgmt:10.131.44.1:443,10.131.44.2:443,10.131.44.3:443|nguid:d22572da-a225-4578-ab1a-9318ac5155c3|proj:cd4eac58-46a5-4a31-b59f-2ec207baa817|scheme:grpcs",
							},
						},
					},
				},
			},
			want: &driver.GetVolumeIDsResponse{
				VolumeIDs: []string{"d22572da-a225-4578-ab1a-9318ac5155c3"},
			},
		},
		{
			name: "invalid lightbits volume",
			req: &driver.GetVolumeIDsRequest{
				PVSpecs: []*corev1.PersistentVolumeSpec{
					{

						PersistentVolumeSource: corev1.PersistentVolumeSource{
							CSI: &corev1.CSIPersistentVolumeSource{
								Driver:       "csi.lightbitslabs.com",
								VolumeHandle: "mgmt:10.131.44.1:443,10.131.44.2:443,10.131.44.3:443|proj:cd4eac58-46a5-4a31-b59f-2ec207baa817|scheme:grpcs",
							},
						},
					},
				},
			},
			want: &driver.GetVolumeIDsResponse{
				VolumeIDs: []string{"mgmt:10.131.44.1:443,10.131.44.2:443,10.131.44.3:443|proj:cd4eac58-46a5-4a31-b59f-2ec207baa817|scheme:grpcs"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{}

			got, err := p.GetVolumeIDs(context.Background(), tt.req)

			if diff := cmp.Diff(tt.wantErr, err); diff != "" {
				t.Errorf("err diff = %s", diff)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("diff = %s", diff)
			}
		})
	}
}
