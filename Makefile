GO=env GO111MODULE=on go
GONOMOD=env GO111MODULE=off go

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: build-cli
build-cli:
	$(GO) build ./dnslink

.PHONY: install
install:
	$(GO) install ./dnslink