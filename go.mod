module github.com/presslabs/stack

go 1.16

require (
	github.com/drone/go-scm v0.0.0
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/karrick/godirwalk v1.15.6 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/pkg/errors v0.9.1
	github.com/presslabs/controller-util v0.3.0-alpha.2
	github.com/presslabs/wordpress-operator v0.11.0-alpha.0
	github.com/prometheus/client_golang v1.9.0
	github.com/rogpeppe/go-internal v1.6.0 // indirect
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	sigs.k8s.io/controller-runtime v0.8.2
)

replace (
	github.com/drone/go-scm v0.0.0 => github.com/presslabs/go-scm v1.5.1-0.20200708152012-713e9c5029bc
	gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7
)
