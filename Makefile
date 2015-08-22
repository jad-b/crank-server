#!/bin/make
BUILD_DIR=build
GOFLAGS=

APPS = cli rest

# Source files per binary. Most require all other Go files outside of bin/
CLI_SRCS = $(wildcard bin/cli/main.go */*.go)
REST_SRCS = $(wildcard bin/rest/main.go */*.go)

# Build all binaries
build: $(APPS)

# How to *actually* build each binary
$(BUILD_DIR)/%:
		@mkdir -p $(dir $@) # Create "$BUILD_DIR/%" if it doesn't exist
		go build ${GOFLAGS} -o $(abspath $(@D))/torque_$(@F) ./$*


# Create a rule for each listed app, which we'll define the targets of below
$(APPS): %: $(BUILD_DIR)/bin/%
$(BINARIES): %: $(BUILD_DIR)/%

# Rules for building each binary
$(BUILD_DIR)/cli: $(CLI_SRCS)
$(BUILD_DIR)/rest: $(REST_SRCS)

clean:
	# Delete built binaries
	rm -rf $(BUILD_DIR)

swagger:
	# Generate full Swagger API spec output file
	./swaggregate.py -m main.yaml -o swagger.yaml

.PHONY: clean all
.PHONY: $(BINARIES)
.PHONY: $(APPS)
