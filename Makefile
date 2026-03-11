default: build

build:
	go build -o terraform-provider-akmatori

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/akmatori/akmatori/0.1.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-akmatori ~/.terraform.d/plugins/registry.terraform.io/akmatori/akmatori/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/terraform-provider-akmatori

test:
	go test ./... -v -timeout 10m

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

docs:
	go generate ./...

generate:
	go generate ./...

fmt:
	gofmt -s -w .

lint:
	golangci-lint run ./...

.PHONY: build install test testacc docs generate fmt lint
