GO=env GO111MODULE=on go

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: install
install:
	$(GO) install ./dnslink