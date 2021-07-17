module github.com/metal-stack/machine-controller-manager-provider-metal

go 1.16

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/gardener/machine-controller-manager v0.40.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/lestrrat-go/jwx v1.0.6 // indirect
	github.com/metal-stack/metal-go v0.15.1
	github.com/metal-stack/metal-lib v0.8.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/prometheus/common v0.15.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
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
)
