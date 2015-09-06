#!/bin/make
BUILD_DIR=build
GOFLAGS=

TORQUE_PKG=github.com/jad-b/torque
TORQUE_APPS=cli rest
APPS=$(addprefix torque_, $(TORQUE_APPS))

# Build all binaries
build: $(addprefix $(BUILD_DIR)/, $(APPS))

# Install all binaries:
install: $(addprefix $(INSTALL_DIR)/, $(APPS))

# How to build each binary
$(BUILD_DIR)/%:
		@mkdir -p $(dir $@) # Create "$BUILD_DIR/%" if it doesn't exist
		go build $(GOFLAGS) -o $(@) ./bin/$*

# How to install each binary
$(INSTALL_DIR)/%:
	go install $(GOFLAGS) $(TORQUE_PKG)/bin/$(@F)


clean:
	# Delete built binaries
	rm -rf $(BUILD_DIR)

swagger:
	# Generate full Swagger API spec output file
	./swaggregate.py -m main.yaml -o swagger.yaml

.PHONY: clean all
.PHONY: $(BINARIES)
.PHONY: $(APPS)
