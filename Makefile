.PHONY: all build run dev dev-server dev-web install test clean gen-http-new gen-http-update hz-install

all: build

build: build-web build-go

build-go:
	go build -o bin/scaffold ./cmd/scaffold

build-server:
	go build -o bin/scaffold-server ./cmd/server

build-web:
	cd web && npm install && npm run build

run: build
	./bin/scaffold --serve

run-server: build-server
	./bin/scaffold-server

dev:
	@echo "Starting development server..."
	@go run ./cmd/server &

dev-cli:
	@go run ./cmd/scaffold --serve --access-key test-key-123

dev-web:
	cd web && npm run dev

install:
	go mod download
	cd web && npm install

test:
	go test ./...

clean:
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/node_modules/
	rm -rf cache/

# Hz code generation
gen-http-new:
	@if [ -z "$(IDL)" ]; then echo "Usage: make gen-http-new IDL=scaffold.proto"; exit 1; fi
	./scripts/gen.sh hz-new $(IDL)

gen-http-update:
	@if [ -z "$(IDL)" ]; then echo "Usage: make gen-http-update IDL=scaffold.proto"; exit 1; fi
	./scripts/gen.sh hz-update $(IDL)

hz-install:
	go install github.com/cloudwego/hertz/cmd/hz@latest
