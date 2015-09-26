#!/bin/bash -eux

FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
TORQUE_PID=


run_server() {
    make install || true
    torque_rest &
    TORQUE_PID=$!
}

poll_server(){
    run_server
    while inotifywait -qre close_write --format "$FORMAT" .; do
        if [ -z ${TORQUE_PID+x} ]; then
            kill $TORQUE_PID
            unset TORQUE_PID
        fi
        run_server
    done
}

poll_server
