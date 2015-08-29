#!/bin/make
BUILD_DIR=build
# ?= sets a default that's overriden by an env var
GOFLAGS?=
TEST_ARGS?=

TORQUE_PKG=github.com/jad-b/torque
TORQUE_APPS=cli rest
APPS=$(addprefix torque_, $(TORQUE_APPS))

# Build all binaries
build: $(addprefix $(BUILD_DIR)/, $(APPS))

# Install all binaries:
install: $(addprefix $(INSTALL_DIR)/, $(APPS))

# How to build each binary
# $@: The target's filename; path included
# $(@F): The filename; path not included
$(BUILD_DIR)/%:
		@mkdir -p $(dir $@) # Create "$BUILD_DIR/%" if it doesn't exist
		go build $(GOFLAGS) -o $@ $(TORQUE_PKG)/bin/$(@F)

# How to install each binary
# $(@F): The filename; path not included
$(INSTALL_DIR)/%:
	go install $(GOFLAGS) $(TORQUE_PKG)/bin/$(@F)

clean:
	# Delete built binaries
	rm -rf $(BUILD_DIR)

redteam:
	go test -v github.com/jad-b/torque/redteam $(TEST_ARGS)

swagger:
	# Generate full Swagger API spec output file
	./swaggregate.py -m main.yaml -o swagger.yaml

.PHONY: clean redteam all
.PHONY: $(BINARIES)
.PHONY: $(APPS)
