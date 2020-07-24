module github.com/presslabs/stack

go 1.14

require (
	cloud.google.com/go v0.39.0 // indirect
	github.com/drone/go-scm v0.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/karrick/godirwalk v1.15.6 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/presslabs/controller-util v0.2.4
	github.com/presslabs/wordpress-operator v0.10.0
	github.com/prometheus/client_golang v1.7.1
	github.com/rogpeppe/go-internal v1.6.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/sys v0.0.0-20200724161237-0e2f3a69832c // indirect
	golang.org/x/tools v0.0.0-20200713235242-6acd2ab80ede // indirect
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/drone/go-scm v0.0.0 => github.com/presslabs/go-scm v1.5.1-0.20200708152012-713e9c5029bc
	gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7
)
