.EXPORT_ALL_VARIABLES:

default: testacc

FILEDATA_BASE_PATH ?= /tmp

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
