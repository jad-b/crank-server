#!/bin/bash -eu
#
# tdd.sh
# ======
# Runs the current tests for the current feature
#
# Because I get tired typing out 'go test -v -tags ... -run ... ...'
FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
TORQUE_PKG="github.com/jad-b/torque"

tdd() {
    go test -v -run TestPostingBodyweight -tags metrics "$TORQUE_PKG/redteam"
}

poll_test(){
    while inotifywait -qre close_write --format "$FORMAT" .; do
        tdd || true
    done
}

tdd
