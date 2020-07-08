BINDIR           ?= $(CURDIR)/../bin

# Run tests
test:
	$(BINDIR)/ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race \
		./pkg/... ./cmd/...

# Build git-webhook binary
build: fmt vet
	go build -o bin/git-webhook .

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

lint:
	GOGC=off $(BINDIR)/golangci-lint run ./pkg/... 
