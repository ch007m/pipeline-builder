GOCMD?=go
TOOL_BIN?=pipe-builder

build: out
	@echo "=====> Building..."
	$(GOCMD) build -ldflags "-s -w -X 'github.com/buildpacks/pack.Version=${PACK_VERSION}' -extldflags '${LDFLAGS}'" -trimpath -o ./out/$(TOOL_BIN) -a .