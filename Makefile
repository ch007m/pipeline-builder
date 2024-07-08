GOCMD?=go
TOOL_BIN?=pipe-builder

build: out
	@echo "=====> Building..."
	$(GOCMD) build -ldflags='-s -w' -trimpath -o ./out/$(TOOL_BIN) -a ./cmd/main.go