.PHONY: all
all: build fmt vet lint test

APP=gohttpclient
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell  go list ./... | grep -v "vendor")
APP_EXECUTABLE="./out/$(APP)"

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

build-deps:
	dep ensure

update-deps:
	dep ensure

compile:
	go build

build: build-deps compile fmt vet lint test

install:
	go install ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@for p in $(UNIT_TEST_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

test:
	ENVIRONMENT=test go test $(UNIT_TEST_PACKAGES) -p=1

test-cover-html:
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(ALL_PACKAGES),\
	ENVIRONMENT=test go test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o out/coverage.html
