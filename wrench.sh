#!/bin/bash -ex

setup(){
    # Install Glide
    GLIDE_TAR='https://github.com/Masterminds/glide/releases/download/0.7.0/glide-0.7.0-linux-amd64.tar.gz'
    mkdir -p /tmp/glide
    curl -Ls $GLIDE_TAR | tar zxv --strip-components 1 -C /tmp/glide

    # Install vendor deps
    GLIDE_HOME=$HOME /tmp/glide/glide up
}

"$@"
