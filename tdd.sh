#!/bin/bash -eu
#
# tdd.sh
# ======
# Runs the current tests for the current feature
#
# Because I get tired typing out 'go test -v -tags ... -run ... ...'
TORQUE_PKG="github.com/jad-b/torque"

tdd() {
    go test -v "$TORQUE_PKG/workouts" -run TestWorkoutCreate
}

poll(){
    local FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
    while inotifywait -qre close_write --format "$FORMAT" .; do
        eval "$@" || true
    done
}
