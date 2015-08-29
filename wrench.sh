#!/bin/bash
#
# wrench.sh provides some helpful shell functions.

setup_test_env(){
    # Start-up the database container
    docker-compose -f deploy/docker-compose.yml up -d db
    # Background the webserver
    make install && torque_rest
}
