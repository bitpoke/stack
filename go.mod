module github.com/bitpoke/stack

go 1.16

require (
	github.com/bitpoke/wordpress-operator v0.12.0
	github.com/drone/go-scm v0.0.0
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/pkg/errors v0.9.1
	github.com/presslabs/controller-util v0.3.0
	github.com/prometheus/client_golang v1.11.1
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	sigs.k8s.io/controller-runtime v0.9.7
)

replace github.com/drone/go-scm v0.0.0 => github.com/bitpoke/go-scm v1.5.1-0.20200708152012-713e9c5029bc
