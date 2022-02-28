KUBECTL=kubectl
KIND=kind
HELM=helm
GO=go
DOCKER=docker
KWCTL=kwctl
VERSION=latest
MANIFEST_TYPE=ClusterAdmissionPolicy
GO_FILES := $(shell find . -type f -name '*.go')

policy.wasm: $(GO_FILES) go.mod go.sum
	$(DOCKER) run --rm -v ${PWD}:/src -w /src tinygo/tinygo:0.18.0 tinygo build \
		-o policy.wasm -target=wasi -no-debug .

annotated-policy.wasm: policy.wasm metadata.yaml
	$(KWCTL) annotate -m metadata.yaml -o annotated-policy.wasm policy.wasm

.PHONY: push-policy
push-policy: annotated-policy.wasm
	$(KWCTL) push annotated-policy.wasm registry://ghcr.io/frelon/kw-palindrome:$(VERSION)
	$(KWCTL) pull registry://ghcr.io/frelon/kw-palindrome:$(VERSION)

manifest.yaml: annotated-policy.wasm
	$(KWCTL) manifest registry://ghcr.io/frelon/kw-palindrome:$(VERSION) --type $(MANIFEST_TYPE) > manifest.yaml

.PHONY: deploy-policy
deploy-policy: manifest.yaml
	$(KUBECTL) apply -f manifest.yaml 

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
