PACKER_SDC_VERSION=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f 2)

.PHONY: build
build:
	@go build

.PHONY: install-packer-sdc
install-packer-sdc:
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${PACKER_SDC_VERSION}

.PHONY: generate
generate: install-packer-sdc
	@go generate ./...
