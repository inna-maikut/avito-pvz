SHELL:=/bin/bash

oapi-codegen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./api/oapi-codegen.yaml ./api/schema.yaml

run-local:
	go run ./cmd/server/main.go

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

lint-ci:
	golangci-lint run ./... --out-format=github-actions --timeout=5m

generate:
	go generate ./...

test-cover-no-integration:
	go test -cover ./...

test-cover:
	go test --tags integration -cover ./...

test-api:
	go test --tags integration ./test/integration

test-repository:
	go test --tags integration ./internal/repository

test-total-cover-no-integration:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out && rm cover.out

test-total-cover:
	go test --tags integration ./... -coverprofile cover.out && go tool cover -func cover.out && rm cover.out

tidy:
	go mod tidy

make_jwt_keys:
	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
	echo "JWT_SECRET=\"`sed -E 's/\$$/\\\n/g' ecprivatekey.pem`\"" >> .env.override
	rm ecprivatekey.pem
