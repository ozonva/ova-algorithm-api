GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on
export GOPROXY=https://proxy.golang.org|direct

PGV_VERSION:="v0.6.1"
GOOGLEAPIS_VERSION="master"
BUF_VERSION:="v0.51.0"

all: generate build

LOCAL_BIN:=$(CURDIR)/bin

build:
	go build -o bin/main cmd/ova-algorithm-api/main.go

run:
	go run cmd/ova-algorithm-api/main.go

.PHONY: generate

generate:
	$(LOCAL_BIN)/mockgen -source=internal/repo/repo.go > internal/mock_repo/mock_repo.go
	$(LOCAL_BIN)/mockgen -source=internal/flusher/flusher.go > internal/mock_flusher/mock_flusher.go
	GOBIN=$(LOCAL_BIN) protoc -I api -I vendor.protogen \
	      --go_out=pkg --go_opt=paths=source_relative  \
	      --go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
	      ova-algorithm-api/ova-algorithm-api.proto
test:
	go test -race ./...

.PHONY:	deps

deps: .install-go-deps

.PHONY: .install-go-deps

.install-go-deps:
	ls go.mod || go mod init github.com/ozonva/ova-algorithm-api
	GOBIN=$(LOCAL_BIN) go get -d google.golang.org/protobuf
	GOBIN=$(LOCAL_BIN) go get -d google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go get -d google.golang.org/grpc
	GOBIN=$(LOCAL_BIN) go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go get -d github.com/rs/zerolog/log
	GOBIN=$(LOCAL_BIN) go get -d github.com/golang/mock/mockgen
	GOBIN=$(LOCAL_BIN) go get -d github.com/pressly/goose/v3/cmd/goose
	GOBIN=$(LOCAL_BIN) go get -d github.com/jackc/pgx/stdlib
	GOBIN=$(LOCAL_BIN) go get -d github.com/onsi/ginkgo/ginkgo
	GOBIN=$(LOCAL_BIN) go get -d github.com/onsi/gomega
	GOBIN=$(LOCAL_BIN) go get -d github.com/Masterminds/squirrel
	GOBIN=$(LOCAL_BIN) go get -d github.com/opentracing/opentracing-go
	GOBIN=$(LOCAL_BIN) go get -d github.com/uber/jaeger-client-go
	GOBIN=$(LOCAL_BIN) go get -d github.com/prometheus/client_golang/prometheus
	GOBIN=$(LOCAL_BIN) go get -d github.com/Shopify/sarama
	GOBIN=$(LOCAL_BIN) go get -d github.com/fsnotify/fsnotify
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose
	GOBIN=$(LOCAL_BIN) go install github.com/onsi/ginkgo/ginkgo
