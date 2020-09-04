package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	api "github.com/metal-stack/machine-controller-manager-provider-metal/pkg/provider/apis"
	"github.com/metal-stack/machine-controller-manager-provider-metal/pkg/provider/apis/validation"
	corev1 "k8s.io/api/core/v1"
)

// decodeProviderSpecAndSecret converts request parameters to api.ProviderSpec & api.Secrets
func decodeProviderSpecAndSecret(machineClass *v1alpha1.MachineClass, secret *corev1.Secret) (*api.MetalProviderSpec, error) {
	var (
		providerSpec *api.MetalProviderSpec
	)

	// Extract providerSpec
	err := json.Unmarshal(machineClass.ProviderSpec.Raw, &providerSpec)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Validate the Spec and Secrets
	ValidationErr := validation.ValidateMetalProviderSpec(providerSpec, secret)
	if ValidationErr != nil {
		err = fmt.Errorf("Error while validating ProviderSpec %v", ValidationErr)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return providerSpec, nil
}

func encodeMachineID(partition, machineID string) string {
	return fmt.Sprintf("metal:///%s/%s", partition, machineID)
}

func decodeMachineID(id string) string {
	splitProviderID := strings.Split(id, "/")
	return splitProviderID[len(splitProviderID)-1]
}
