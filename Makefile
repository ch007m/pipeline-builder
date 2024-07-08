GOCMD?=go
TOOL_BIN?=pipe-builder

build: out
	@echo "=====> Generate..."
	$(GOCMD) generate templates/lifecycle/task/*.go

	@echo "=====> Building..."
	$(GOCMD) build -ldflags='-s -w' -trimpath -o ./out/$(TOOL_BIN) -a ./cmd/main.go