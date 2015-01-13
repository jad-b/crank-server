#!/bin/bash -eux

setup(){
    # Install Glide
	GLIDE_VERSION="0.8.2"
    GLIDE_TAR="https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/glide-$GLIDE_VERSION-linux-amd64.tar.gz"
    mkdir -p /tmp/glide
	# Download glide archive
	curl -Ls -o /tmp/glide/glide-$GLIDE_VESION.tar.gz $GLIDE_TAR
	# Extract glide
	tar zxv --strip-components 1 -C /tmp/glide -f /tmp/glide/glide-$GLIDE_VERSION.tar.gz

    # Install vendor deps
    GLIDE_HOME=$HOME /tmp/glide/glide up
}

"$@"
