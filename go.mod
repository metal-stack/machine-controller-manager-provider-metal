module github.com/metal-stack/machine-controller-manager-provider-metal

go 1.15

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/gardener/machine-controller-manager v0.35.2
	github.com/go-openapi/validate v0.20.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/lestrrat-go/jwx v1.0.6 // indirect
	github.com/metal-stack/masterdata-api v0.8.4 // indirect
	github.com/metal-stack/metal-go v0.11.1
	github.com/metal-stack/metal-lib v0.6.6
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	golang.org/x/sys v0.0.0-20210104204734-6f8348627aad // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/api v0.17.14
	k8s.io/apimachinery v0.17.14
	k8s.io/client-go v1.5.1 // indirect
	k8s.io/cluster-bootstrap v0.17.14 // indirect
	k8s.io/component-base v0.17.14
	k8s.io/klog v1.0.0
	sigs.k8s.io/yaml v1.2.0 // indirect
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
