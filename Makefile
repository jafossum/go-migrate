TEST_PKGS=./pkg/...
GEN_PKGS=./pkg/...
VET_PKGS=./pkg/... ./cmd/...

.PHONY: all test fmt docker

all: fmt vet test 

test:
	go test -v $(TEST_PKGS)

vet:
	go vet $(VET_PKGS)

fmt:
	go fmt $(VET_PKGS)
