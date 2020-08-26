// Package validation - validation is used to validate cloud specific ProviderSpec
package validation

import (
	"fmt"

	api "github.com/metal-stack/machine-controller-manager-provider-metal/pkg/provider/apis"
	corev1 "k8s.io/api/core/v1"
)

// ValidateMetalProviderSpec validates provider spec and secret to check if all fields are present and valid
func ValidateMetalProviderSpec(spec *api.MetalProviderSpec, secrets *corev1.Secret) []error {
	var allErrs []error

	if "" == spec.Image {
		allErrs = append(allErrs, fmt.Errorf("image is required field"))
	}
	if "" == spec.Network {
		allErrs = append(allErrs, fmt.Errorf("network is required field"))
	}
	if "" == spec.Partition {
		allErrs = append(allErrs, fmt.Errorf("partition is required field"))
	}
	if "" == spec.Project {
		allErrs = append(allErrs, fmt.Errorf("project is required field"))
	}
	if "" == spec.Size {
		allErrs = append(allErrs, fmt.Errorf("size is required field"))
	}

	allErrs = append(allErrs, validateSecrets(secrets)...)

	return allErrs
}

func validateSecrets(secret *corev1.Secret) []error {
	var allErrs []error
	if ("" == string(secret.Data["metalAPIHMac"])) == ("" == string(secret.Data["metalAPIKey"])) {
		allErrs = append(allErrs, fmt.Errorf("Either metalAPIHMac or metalAPIKey is required field"))
	}
	if "" == string(secret.Data["metalAPIURL"]) {
		allErrs = append(allErrs, fmt.Errorf("Secret metalAPIURL is required field"))
	}
	return allErrs
}
