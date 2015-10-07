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
   torque_cli -v -web -username jdb -password torqued \
        update bodyweight -timestamp '2015-10-04 12:22:05' \
           -weight 177.8 \
           -comment "Have some U with your CR"
}

poll_test(){
    while inotifywait -qre close_write --format "$FORMAT" .; do
        tdd || true
    done
}

tdd
