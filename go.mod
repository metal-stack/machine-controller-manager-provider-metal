module github.com/metal-stack/machine-controller-manager-provider-metal

go 1.16

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/gardener/machine-controller-manager v0.39.0
	github.com/go-openapi/validate v0.20.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/jwx v1.0.6 // indirect
	github.com/metal-stack/metal-go v0.14.0
	github.com/metal-stack/metal-lib v0.8.0
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/prometheus/common v0.25.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/api v0.17.17
	k8s.io/apimachinery v0.17.17
	k8s.io/client-go v1.5.2 // indirect
	k8s.io/cluster-bootstrap v0.17.17 // indirect
	k8s.io/component-base v0.17.17
	k8s.io/klog v1.0.0
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.17
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.17
	k8s.io/apiserver => k8s.io/apiserver v0.17.17
	k8s.io/client-go => k8s.io/client-go v0.17.17
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.17
	k8s.io/code-generator => k8s.io/code-generator v0.17.17
)
