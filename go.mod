module github.com/presslabs/stack

go 1.14

require (
	cloud.google.com/go v0.39.0 // indirect
	github.com/appscode/jsonpatch v0.0.0-20190108182946-7c0e3b262f30 // indirect
	github.com/drone/go-scm v0.0.0
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gobuffalo/genny v0.0.0-20190124191459-3310289fa4b4 // indirect
	github.com/gobuffalo/meta v0.0.0-20190121163014-ecaa953cbfb3 // indirect
	github.com/gobuffalo/packr/v2 v2.0.0-rc.15
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.5.0
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/presslabs/controller-util v0.1.13
	github.com/presslabs/wordpress-operator v0.3.4
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20190115171406-56726106282f // indirect
	github.com/prometheus/common v0.1.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190117184657-bf6a532e95b1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	golang.org/x/oauth2 v0.0.0-20190523182746-aaccbc9213b0 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/appengine v1.6.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20181213150558-05914d821849 // indirect
	k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476 // indirect
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	k8s.io/klog v0.3.1 // indirect
	k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22 // indirect
	sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/testing_frameworks v0.1.1 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	github.com/drone/go-scm v0.0.0 => github.com/presslabs/go-scm v1.5.1-0.20200708152012-713e9c5029bc
	gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7
)
