module github.com/presslabs/stack

go 1.14

require (
	cloud.google.com/go v0.39.0 // indirect
	github.com/drone/go-scm v0.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/presslabs/controller-util v0.2.4
	github.com/presslabs/wordpress-operator v0.10.0
	github.com/prometheus/client_golang v1.7.1
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/drone/go-scm v0.0.0 => github.com/presslabs/go-scm v1.5.1-0.20200708152012-713e9c5029bc
	gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7
)
