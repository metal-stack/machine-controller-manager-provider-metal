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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/metal-stack/metal-lib/pkg/tag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"

	api "github.com/metal-stack/machine-controller-manager-provider-metal/pkg/provider/apis"
	metalv1alpha1 "github.com/metal-stack/machine-controller-manager-provider-metal/pkg/provider/migration/legacy-api/machine/v1alpha1"
)

// NOTE
//
// The basic working of the controller will work with just implementing the CreateMachine() & DeleteMachine() methods.
// You can first implement these two methods and check the working of the controller.
// Leaving the other methods to NOT_IMPLEMENTED error status.
// Once this works you can implement the rest of the methods.
//
// Also make sure each method return appropriate errors mentioned in `https://github.com/gardener/machine-controller-manager/blob/master/docs/development/machine_error_codes.md`

// CreateMachine handles a machine creation request
// REQUIRED METHOD
//
// REQUEST PARAMETERS (driver.CreateMachineRequest)
// Machine               *v1alpha1.Machine        Machine object from whom VM is to be created
// MachineClass          *v1alpha1.MachineClass   MachineClass backing the machine object
// Secret                *corev1.Secret           Kubernetes secret that contains any sensitive data/credentials
//
// RESPONSE PARAMETERS (driver.CreateMachineResponse)
// ProviderID            string                   Unique identification of the VM at the cloud provider. This could be the same/different from req.MachineName.
//                                                ProviderID typically matches with the node.Spec.ProviderID on the node object.
//                                                Eg: gce://project-name/region/vm-ProviderID
// NodeName              string                   Returns the name of the node-object that the VM register's with Kubernetes.
//                                                This could be different from req.MachineName as well
// LastKnownState        string                   (Optional) Last known state of VM during the current operation.
//                                                Could be helpful to continue operations in future requests.
//
// OPTIONAL IMPLEMENTATION LOGIC
// It is optionally expected by the safety controller to use an identification mechanisms to map the VM Created by a providerSpec.
// These could be done using tag(s)/resource-groups etc.
// This logic is used by safety controller to delete orphan VMs which are not backed by any machine CRD
//
func (p *Provider) CreateMachine(ctx context.Context, req *driver.CreateMachineRequest) (*driver.CreateMachineResponse, error) {
	klog.V(2).Infof("machine creation request has been recieved for %q", req.Machine.Name)
	providerSpec, err := decodeProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	m, err := p.initDriver(req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	networks := []metalgo.MachineAllocationNetwork{
		{
			Autoacquire: true,
			NetworkID:   providerSpec.Network,
		},
	}
	userData := strings.TrimSpace(string(req.Secret.Data["userData"]))

	createRequest := &metalgo.MachineCreateRequest{
		Description:   req.Machine.Name + " created by Gardener.",
		Name:          req.Machine.Name,
		Hostname:      req.Machine.Name,
		UserData:      userData,
		Size:          providerSpec.Size,
		Project:       providerSpec.Project,
		Networks:      networks,
		Partition:     providerSpec.Partition,
		Image:         providerSpec.Image,
		Tags:          providerSpec.Tags,
		SSHPublicKeys: providerSpec.SSHKeys,
	}

	mcr, err := m.MachineCreate(createRequest)
	if err != nil {
		klog.Errorf("could not create machine: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	klog.V(2).Infof("machine creation request has been processed for %q", req.Machine.Name)

	return &driver.CreateMachineResponse{
		ProviderID: encodeMachineID(providerSpec.Partition, *mcr.Machine.ID),
		NodeName:   *mcr.Machine.Allocation.Name,
	}, nil
}

// DeleteMachine handles a machine deletion request
//
// REQUEST PARAMETERS (driver.DeleteMachineRequest)
// Machine               *v1alpha1.Machine        Machine object from whom VM is to be deleted
// MachineClass          *v1alpha1.MachineClass   MachineClass backing the machine object
// Secret                *corev1.Secret           Kubernetes secret that contains any sensitive data/credentials
//
// RESPONSE PARAMETERS (driver.DeleteMachineResponse)
// LastKnownState        bytes(blob)              (Optional) Last known state of VM during the current operation.
//                                                Could be helpful to continue operations in future requests.
//
func (p *Provider) DeleteMachine(ctx context.Context, req *driver.DeleteMachineRequest) (*driver.DeleteMachineResponse, error) {
	klog.V(2).Infof("machine deletion request has been recieved for %q", req.Machine.Name)
	providerSpec, err := decodeProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	clusterIDTag, err := extractClusterTag(providerSpec.Tags)
	if err != nil {
		klog.V(2).Infof("machine deletion request for machine %q failed because provider spec did not contain metal-stack cluster tag", req.Machine.Name)
		return nil, status.Error(codes.Internal, "machine deletion request failed because provider spec did not contain metal-stack cluster tag for")
	}

	m, err := p.initDriver(req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	id := decodeMachineID(req.Machine.Spec.ProviderID)

	mfr := &metalgo.MachineFindRequest{
		ID:                &id,
		AllocationProject: &providerSpec.Project,
		Tags:              []string{clusterIDTag},
	}

	resp, err := m.MachineFind(mfr)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	switch len(resp.Machines) {
	case 0:
		klog.Infof("no machine with id %q found in project %q, already deleted and therefore skipping deletion", id, providerSpec.Project)
		return &driver.DeleteMachineResponse{}, nil
	case 1:
		_, err = m.MachineDelete(id)
		if err != nil {
			klog.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		klog.Infof("deleted machine %q (%q)", req.Machine.Name, id)
		return &driver.DeleteMachineResponse{}, nil
	default:
		klog.Errorf("error finding machine to delete because more than one search result")
		return nil, status.Error(codes.Internal, "error finding machine to delete because more than one search result")
	}
}

// GetMachineStatus handles a machine get status request
// OPTIONAL METHOD
//
// REQUEST PARAMETERS (driver.GetMachineStatusRequest)
// Machine               *v1alpha1.Machine        Machine object from whom VM status needs to be returned
// MachineClass          *v1alpha1.MachineClass   MachineClass backing the machine object
// Secret                *corev1.Secret           Kubernetes secret that contains any sensitive data/credentials
//
// RESPONSE PARAMETERS (driver.GetMachineStatueResponse)
// ProviderID            string                   Unique identification of the VM at the cloud provider. This could be the same/different from req.MachineName.
//                                                ProviderID typically matches with the node.Spec.ProviderID on the node object.
//                                                Eg: gce://project-name/region/vm-ProviderID
// NodeName             string                    Returns the name of the node-object that the VM register's with Kubernetes.
//                                                This could be different from req.MachineName as well
//
// The request should return a NOT_FOUND (5) status error code if the machine is not existing
func (p *Provider) GetMachineStatus(ctx context.Context, req *driver.GetMachineStatusRequest) (*driver.GetMachineStatusResponse, error) {
	klog.V(2).Infof("get request has been recieved for %q", req.Machine.Name)
	providerSpec, err := decodeProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	clusterIDTag, err := extractClusterTag(providerSpec.Tags)
	if err != nil {
		klog.V(2).Infof("get request for machine %q failed because provider spec did not contain metal-stack cluster tag", req.Machine.Name)
		return nil, status.Error(codes.Internal, "get machine request failed because provider spec did not contain metal-stack cluster tag for")
	}

	m, err := p.initDriver(req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	id := decodeMachineID(req.Machine.Spec.ProviderID)

	if id == "" {
		return nil, status.Error(codes.NotFound, "machine not found, not yet created")
	}

	resp, err := m.MachineGet(id)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.Machine.Allocation == nil {
		klog.V(2).Infof("machine already released: %q", req.Machine.Name)
		return nil, status.Error(codes.NotFound, "machine already released")
	}

	machineClusterIDTag, err := extractClusterTag(resp.Machine.Tags)
	if err != nil {
		klog.V(2).Infof("machine has no cluster tag anymore: %q", req.Machine.Name)
		return nil, status.Error(codes.NotFound, "machine has no cluster tag anymore")
	}

	if machineClusterIDTag != clusterIDTag {
		klog.V(2).Infof("machine does not belong to this cluster anymore: %q", req.Machine.Name)
		return nil, status.Error(codes.NotFound, "machine does not belong to this cluster anymore")
	}

	klog.V(2).Infof("machine get request has been processed successfully for %q", req.Machine.Name)

	return &driver.GetMachineStatusResponse{
		ProviderID: encodeMachineID(*resp.Machine.Partition.ID, *resp.Machine.ID),
		NodeName:   *resp.Machine.Allocation.Name,
	}, nil
}

// ListMachines lists all the machines possibilly created by a providerSpec
// Identifying machines created by a given providerSpec depends on the OPTIONAL IMPLEMENTATION LOGIC
// you have used to identify machines created by a providerSpec. It could be tags/resource-groups etc
// OPTIONAL METHOD
//
// REQUEST PARAMETERS (driver.ListMachinesRequest)
// MachineClass          *v1alpha1.MachineClass   MachineClass based on which VMs created have to be listed
// Secret                *corev1.Secret           Kubernetes secret that contains any sensitive data/credentials
//
// RESPONSE PARAMETERS (driver.ListMachinesResponse)
// MachineList           map<string,string>  A map containing the keys as the MachineID and value as the MachineName
//                                           for all machine's who where possibilly created by this ProviderSpec
//
func (p *Provider) ListMachines(ctx context.Context, req *driver.ListMachinesRequest) (*driver.ListMachinesResponse, error) {
	klog.V(2).Infof("list machines request has been recieved for %q", req.MachineClass.Name)
	providerSpec, err := decodeProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	m, err := p.initDriver(req.Secret)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	listOfVMs := make(map[string]string)

	clusterIDTag, err := extractClusterTag(providerSpec.Tags)
	if err != nil {
		klog.V(2).Infof("list machines request failed because provider spec did not contain metal-stack cluster tag for %q", req.MachineClass.Name)
		return nil, status.Error(codes.Internal, "list machines request failed because provider spec did not contain metal-stack cluster tag for")
	}

	findRequest := &metalgo.MachineFindRequest{
		AllocationProject: &providerSpec.Project,
		Tags:              []string{clusterIDTag},
	}
	resp, err := m.MachineFind(findRequest)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, m := range resp.Machines {
		if !stringSliceContains(m.Tags, "kubernetes.io/role=node") {
			continue
		}

		if m.ID == nil || m.Partition == nil || m.Partition.ID == nil {
			continue
		}

		providerID := encodeMachineID(*m.Partition.ID, *m.ID)
		listOfVMs[providerID] = *m.Allocation.Hostname
	}

	klog.V(2).Infof("list machines request has been processed successfully for %q, found %v", req.MachineClass.Name, listOfVMs)

	return &driver.ListMachinesResponse{MachineList: listOfVMs}, nil
}

func stringSliceContains(s []string, val string) bool {
	for _, item := range s {
		if item == val {
			return true
		}
	}
	return false
}

// GetVolumeIDs returns a list of Volume IDs for all PV Specs for whom an provider volume was found
//
// REQUEST PARAMETERS (driver.GetVolumeIDsRequest)
// PVSpecList            []*corev1.PersistentVolumeSpec       PVSpecsList is a list PV specs for whom volume-IDs are required.
//
// RESPONSE PARAMETERS (driver.GetVolumeIDsResponse)
// VolumeIDs             []string                             VolumeIDs is a repeated list of VolumeIDs.
//
func (p *Provider) GetVolumeIDs(ctx context.Context, req *driver.GetVolumeIDsRequest) (*driver.GetVolumeIDsResponse, error) {
	// Log messages to track start and end of request
	klog.V(2).Infof("GetVolumeIDs request has been recieved for %q", req.PVSpecs)
	volumeIDs := []string{}
	specs := req.PVSpecs
	for i := range specs {
		spec := specs[i]
		if spec.CSI == nil {
			// Not a CSI volume
			continue
		}

		volumeIDs = append(volumeIDs, spec.CSI.VolumeHandle)
	}
	klog.V(2).Infof("GetVolumeIDs request has been processed successfully for %q", req.PVSpecs)
	return &driver.GetVolumeIDsResponse{VolumeIDs: volumeIDs}, nil
}

// GenerateMachineClassForMigration helps in migration of one kind of machineClass CR to another kind.
// For instance an machineClass custom resource of `AWSMachineClass` to `MachineClass`.
// Implement this functionality only if something like this is desired in your setup.
// If you don't require this functionality leave is as is. (return Unimplemented)
//
// The following are the tasks typically expected out of this method
// 1. Validate if the incoming classSpec is valid one for migration (e.g. has the right kind).
// 2. Migrate/Copy over all the fields/spec from req.ProviderSpecificMachineClass to req.MachineClass
// For an example refer
//		https://github.com/prashanth26/machine-controller-manager-provider-gcp/blob/migration/pkg/gcp/machine_controller.go#L222-L233
//
// REQUEST PARAMETERS (driver.GenerateMachineClassForMigration)
// ProviderSpecificMachineClass    interface{}                             ProviderSpecificMachineClass is provider specfic machine class object (E.g. AWSMachineClass). Typecasting is required here.
// MachineClass 				   *v1alpha1.MachineClass                  MachineClass is the machine class object that is to be filled up by this method.
// ClassSpec                       *v1alpha1.ClassSpec                     Somemore classSpec details useful while migration.
//
// RESPONSE PARAMETERS (driver.GenerateMachineClassForMigration)
// NONE
//
func (p *Provider) GenerateMachineClassForMigration(ctx context.Context, req *driver.GenerateMachineClassForMigrationRequest) (*driver.GenerateMachineClassForMigrationResponse, error) {
	// Log messages to track start and end of request
	klog.V(2).Infof("MigrateMachineClass request has been recieved for %q", req.ClassSpec)

	metalMachineClass := req.ProviderSpecificMachineClass.(*metalv1alpha1.MetalMachineClass)

	// Check if incoming CR is valid CR for migration
	// In this case, the MachineClassKind to be matching
	if req.ClassSpec.Kind != "MetalMachineClass" {
		err := status.Error(codes.Internal, "Migration cannot be done for this machineClass kind")
		klog.Error(err.Error())
		return nil, err
	}

	providerSpec := &api.MetalProviderSpec{
		Partition: metalMachineClass.Spec.Partition,
		Size:      metalMachineClass.Spec.Size,
		Image:     metalMachineClass.Spec.Image,
		Project:   metalMachineClass.Spec.Project,
		Network:   metalMachineClass.Spec.Network,
		Tags:      metalMachineClass.Spec.Tags,
		SSHKeys:   metalMachineClass.Spec.SSHKeys,
	}

	providerSpecMarshal, err := json.Marshal(providerSpec)
	if err != nil {
		klog.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Migrate finalizers, labels, annotations
	req.MachineClass.Name = metalMachineClass.Name
	req.MachineClass.Labels = metalMachineClass.Labels
	req.MachineClass.Annotations = metalMachineClass.Annotations
	req.MachineClass.Finalizers = metalMachineClass.Finalizers
	req.MachineClass.ProviderSpec = runtime.RawExtension{
		Raw: providerSpecMarshal,
	}
	req.MachineClass.SecretRef = metalMachineClass.Spec.SecretRef
	req.MachineClass.Provider = "metal"

	klog.V(2).Infof("MigrateMachineClass request has been processed successfully for %q", req.ClassSpec)
	return &driver.GenerateMachineClassForMigrationResponse{}, nil
}

func extractClusterTag(tags []string) (string, error) {
	for _, t := range tags {
		if strings.HasPrefix(t, tag.ClusterID) {
			return t, nil
		}
	}
	return "", fmt.Errorf("tag not found")
}
