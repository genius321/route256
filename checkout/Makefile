CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=route256/checkout/cmd/app

all: format build test lint

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin
bin:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
 	cd vendor-proto/googleapis &&\
	git sparse-checkout set --no-cone google/api &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/validate-repo &&\
	cd vendor-proto/validate-repo &&\
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p  vendor-proto
	mv vendor-proto/validate-repo/validate vendor-proto
	rm -rf vendor-proto/validate-repo


generate-loms: bin vendor-proto/google/api vendor-proto/google/protobuf
	mkdir -p internal/pkg/loms
	$(PROTOC) -I api/loms -I vendor-proto \
	--go_out internal/pkg/loms --go_opt paths=source_relative \
	--go-grpc_out internal/pkg/loms --go-grpc_opt paths=source_relative \
	--grpc-gateway_out internal/pkg/loms --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:internal/pkg/loms" \
	api/loms/loms.proto

generate-checkout: bin vendor-proto/google/api vendor-proto/google/protobuf
	mkdir -p internal/pkg/checkout
	$(PROTOC) -I api/checkout -I vendor-proto \
	--go_out internal/pkg/checkout --go_opt paths=source_relative \
	--go-grpc_out internal/pkg/checkout --go-grpc_opt paths=source_relative \
	--grpc-gateway_out internal/pkg/checkout --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:internal/pkg/checkout" \
	api/checkout/checkout.proto

generate-product: bin vendor-proto/google/api vendor-proto/google/protobuf
	mkdir -p internal/pkg/product
	$(PROTOC) -I api/product-service -I vendor-proto \
	--go_out internal/pkg/product --go_opt paths=source_relative \
	--go-grpc_out internal/pkg/product --go-grpc_opt paths=source_relative \
	--grpc-gateway_out internal/pkg/product --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:internal/pkg/product" \
	api/product-service/product-service.proto

generate: generate-loms generate-checkout generate-product