BINDIR           ?= $(CURDIR)/../bin

# Run tests
test:
	$(BINDIR)/ginkgo version
	$(BINDIR)/ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race \
		./pkg/...

# Build git-webhook binary
build: fmt vet
	go build -o bin/git-webhook .

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./pkg/...

# Run go vet against code
vet:
	go vet ./pkg/...

lint:
	$(BINDIR)/golangci-lint version
	GOGC=off $(BINDIR)/golangci-lint run ./pkg/... 
