#!/bin/bash -eu

EXCLUDE_DIRS='.git/**'
FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
TORQUE_PID=


run_server() {
	local APP="$1"
	echo "Running $APP"
    make install || true
    "$APP" &
    TORQUE_PID=$!
    echo "$APP PID: $TORQUE_PID"
}

poll_server(){
	local APP="${1:-torque_rest}"
	echo "Polling on $APP"
    # kill any other server, just in case
    if pgrep -u "$USER $APP"; then
        pkill -u "$USER" --signal SIGKILL "$APP"
    fi
    run_server "$APP"
    # List all dirs to watch
    while inotifywait -qre close_write --exclude "$EXCLUDE_DIRS" --format "$FORMAT" .; do
        if ! [  -z "${TORQUE_PID+x}" ]; then
            echo "Killing torque server..."
            kill "$TORQUE_PID" || true
            unset TORQUE_PID
        fi
        run_server "$APP"
    done
}

poll_server "$@"
