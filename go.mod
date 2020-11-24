module github.com/metal-stack/machine-controller-manager-provider-metal

go 1.15

require (
	github.com/gardener/machine-controller-manager v0.35.0
	github.com/metal-stack/metal-go v0.11.1
	github.com/metal-stack/metal-lib v0.6.6
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.17.14
	k8s.io/apimachinery v0.17.14
	k8s.io/component-base v0.17.14
	k8s.io/klog v1.0.0
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.14
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.14
	k8s.io/apiserver => k8s.io/apiserver v0.17.14
	k8s.io/client-go => k8s.io/client-go v0.17.14
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.14
	k8s.io/code-generator => k8s.io/code-generator v0.17.14
	// k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
)
