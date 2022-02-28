KUBECTL=kubectl
KIND=kind
HELM=helm
GO=go
SOURCE_FILES := $(shell find . -type f -name '*.go')

policy.wasm: $(SOURCE_FILES) go.mod go.sum
	docker run --rm -v ${PWD}:/src -w /src tinygo/tinygo:0.18.0 tinygo build \
		-o policy.wasm -target=wasi -no-debug .

annotated-policy.wasm: policy.wasm metadata.yaml
	kwctl annotate -m metadata.yaml -o annotated-policy.wasm policy.wasm

.PHONY: test
test:
	$(GO) test -race ./...

.PHONY: up
up: cluster cert-manager kubewarden

.PHONY: cluster
cluster:
	$(KIND) create cluster

.PHONY: cert-manager
cert-manager:
	$(KUBECTL) apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
	$(KUBECTL) wait --for=condition=Available deployment --timeout=2m -n cert-manager --all

.PHONY: kubewarden
kubewarden:
	$(HELM) repo add kubewarden https://charts.kubewarden.io
	$(HELM) install --wait -n kubewarden --create-namespace kubewarden-crds kubewarden/kubewarden-crds
	$(HELM) install --wait -n kubewarden kubewarden-controller kubewarden/kubewarden-controller

.PHONY: down
down:
	$(KIND) delete cluster
