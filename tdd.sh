#!/bin/bash -eu
#
# tdd.sh
# ======
# Runs the current tests for the current feature
#
# Because I get tired typing out 'go test -v -tags ... -run ... ...'
TORQUE_PKG="github.com/jad-b/torque"


poll_test(){
    while inotifywait -qre close_write --format "$FORMAT" .; do
        go test -v -run TestPostingBodyweight -tags metrics "$TORQUE_PKG/redteam"
    done
}
