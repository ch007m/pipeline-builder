GOCMD?=go
TOOL_BIN?=pipe-builder

build: check-statik out
	@echo "=====> Generate..."
	$(GOCMD) generate templates/lifecycle/task/*.go

	@echo "=====> Building..."
	$(GOCMD) build -ldflags='-s -w' -trimpath -o ./out/$(TOOL_BIN) -a ./cmd/main.go

## out: Make a directory for output
out:
	@mkdir out || (exit 0)

check-statik:
	@if ! command -v statik &> /dev/null; then \
		echo "=====> statik not found. Installing statik..."; \
		if ! $(GOCMD) install github.com/rakyll/statik@latest; then \
			echo "Error: failed to install statik"; \
			exit 1; \
    	fi \
	else \
		echo "=====> statik (github.com/rakyll/statik) is already installed."; \
	fi

.PHONY: check-statik build out