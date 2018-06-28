.PHONY: all

CURDIR := $(shell pwd)

help: ## Print this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

init: ## Initialize project
	go get -u github.com/golang/dep/cmd/dep
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	cd ./v1 && dep ensure

mockgen: ## Generate mocks
	cd ./v1; \
	mockgen -package datastore -destination osindatastore_mock_test.go go.mercari.io/datastore Client; \
	mockgen -source storage.go -package datastore -destination storage_mock_test.go

test: ## Execute test
	go test ./...
