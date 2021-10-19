module github.com/metal-stack/machine-controller-manager-provider-metal

go 1.16

require (
	github.com/gardener/machine-controller-manager v0.41.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/metal-stack/metal-go v0.15.7
	github.com/metal-stack/metal-lib v0.9.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
	k8s.io/client-go v1.5.1 // indirect
	k8s.io/cluster-bootstrap v0.17.14 // indirect
	k8s.io/component-base v0.20.6
	k8s.io/klog v1.0.0
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	k8s.io/api => k8s.io/api v0.20.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.6
	k8s.io/apiserver => k8s.io/apiserver v0.20.6
	k8s.io/client-go => k8s.io/client-go v0.20.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.6
	k8s.io/code-generator => k8s.io/code-generator v0.20.6
)
