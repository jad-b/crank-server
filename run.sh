#!/bin/bash -eu

EXCLUDE_DIRS='.git'
FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
TORQUE_PID=


run_server() {
    make install || true
    torque_rest &
    TORQUE_PID=$!
    echo "Torque REST server PID: $TORQUE_PID"
}

poll_server(){
    # kill any other server, just in case
    if [ $(pgrep -u $USER torque_rest) ]; then
        kill $(pgrep -u $USER torque_rest)
    fi
    run_server
    # List all dirs to watch
    while inotifywait -qre close_write --format "$FORMAT" .; do
        if ! [  -z ${TORQUE_PID+x} ]; then
            echo "Killing torque server..."
            kill $TORQUE_PID
            unset TORQUE_PID
        fi
        run_server
    done
}

poll_server
