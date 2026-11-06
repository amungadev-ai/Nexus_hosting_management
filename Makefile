# NexusPanel — Go monorepo (walking skeleton)
GO ?= go
BIN := bin

build:
	$(GO) build -o $(BIN)/nexus-control-plane ./cmd/control-plane
	$(GO) build -o $(BIN)/nexus-host-agent ./cmd/host-agent
	$(GO) build -o $(BIN)/nexus-npctl ./cmd/npctl
	$(GO) build -o $(BIN)/nexus-installer ./cmd/installer

test:
	$(GO) test ./... -count=1

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

# Local end-to-end smoke run (SQLite, throwaway state dir).
smoke: build
	@rm -rf tmp/smoke
	./$(BIN)/nexus-control-plane -data-dir tmp/smoke -addr 127.0.0.1:8080 -site-addr 127.0.0.1:8081

clean:
	rm -rf $(BIN) tmp/smoke

.PHONY: build test vet tidy smoke clean
